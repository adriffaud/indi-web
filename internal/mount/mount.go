package mount

import (
	"log/slog"

	"github.com/a-h/templ"
	"github.com/adriffaud/indi-web/internal/coordconv"
	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
	"github.com/adriffaud/indi-web/ui/components"
)

type Mount struct {
	client    *indiclient.Client
	eventChan chan indiclient.Event
	htmlChan  chan templ.Component
	Connected bool
	Driver    string
	Parked    bool
	Parking   bool
	Tracking  bool
	RA        string
	DEC       string
}

func NewMount(driver string, eventChan chan indiclient.Event, htmlChan chan templ.Component) *Mount {
	mount := Mount{
		eventChan: eventChan,
		htmlChan:  htmlChan,
		Driver:    driver,
		RA:        "00:00:00",
		DEC:       "00:00:00",
		Parked:    true,
		Parking:   false,
		Tracking:  false,
		Connected: false,
	}

	go func() {
		for {
			event := <-mount.eventChan
			mount.onEvent(event)
		}
	}()

	return &mount
}

func (m *Mount) onEvent(event indiclient.Event) {
	if event.EventType == indiclient.Timeout && !m.Connected {
		slog.Debug("🤓 Mount not connected, connecting...")
		err := m.client.Connect(m.Driver)
		if err != nil {
			slog.Error("🔴 Could not automatically connect to mount", "err", err)
		}
		m.Connected = true
	}

	if m.Driver == "" || event.Property.Device != m.Driver {
		return
	}

	switch event.Property.Name {
	case "EQUATORIAL_EOD_COORD":
		for _, value := range event.Property.Values {
			formated, err := coordconv.DecimalToSexagesimal(value.Value)
			if err != nil {
				slog.Error("could not convert coords in sexagesimal", "err", err, "name", value.Name, "value", value)
			}

			var component templ.Component

			if value.Name == "RA" {
				m.RA = formated
				component = components.TextInput("ra_input", m.RA, templ.Attributes{"disabled": "true", "hx-swap-oob": "true"})
			}
			if value.Name == "DEC" {
				m.DEC = formated
				component = components.TextInput("dec_input", m.DEC, templ.Attributes{"disabled": "true", "hx-swap-oob": "true"})
			}

			m.htmlChan <- component
		}
	case "TELESCOPE_PARK":
		slog.Debug("⚠️ TELESCOPE_PARK event", "property", event.Property)
		for _, value := range event.Property.Values {
			if value.Value != "On" {
				continue
			}

			if value.Name == "UNPARK" {
				m.Parked = false
			} else if value.Name == "PARK" {
				m.Parked = true
				m.Parking = false
			}

			m.htmlChan <- components.ParkButton(m.Parked, m.Parking)
			m.htmlChan <- components.TrackButton(m.Tracking, m.Parked)
		}
	default:
		slog.Debug("🚧 Unhandled mount event", "name", event.Property.Name)
	}

}

func (m *Mount) SetClient(client *indiclient.Client) {
	m.client = client
}

func (m *Mount) Park() {
	if m.Parked || m.Parking {
		return
	}

	m.Parking = true
	err := m.client.NewPropertyValue(indiclient.PropertySelector{Device: m.Driver, Name: "TELESCOPE_PARK", ValueName: "PARK"})
	if err != nil {
		slog.Error("🔴 could not start mount parking", "error", err)
	}
}

func (m *Mount) Unpark() {
	if !m.Parked {
		return
	}

	err := m.client.NewPropertyValue(indiclient.PropertySelector{Device: m.Driver, Name: "TELESCOPE_PARK", ValueName: "UNPARK"})
	if err != nil {
		slog.Error("🔴 could not start mount unparking", "error", err)
	}
}

func (m *Mount) StartTracking() {
	if m.Tracking || m.Parked || m.Parking {
		return
	}

	err := m.client.NewPropertyValue(indiclient.PropertySelector{Device: m.Driver, Name: "TELESCOPE_TRACK_STATE", ValueName: "TRACK_ON"})
	if err != nil {
		slog.Error("🔴 could not start tracking", "error", err)
	}

	// INDI doesn't update property state for tracking state, so updating state manually.
	m.Tracking = true
	m.htmlChan <- components.TrackButton(m.Tracking, m.Parked)
}

func (m *Mount) StopTracking() {
	if !m.Tracking {
		return
	}

	err := m.client.NewPropertyValue(indiclient.PropertySelector{Device: m.Driver, Name: "TELESCOPE_TRACK_STATE", ValueName: "TRACK_OFF"})
	if err != nil {
		slog.Error("🔴 could not stop tracking", "error", err)
	}

	// INDI doesn't update property state for tracking state, so updating state manually.
	m.Tracking = false
	m.htmlChan <- components.TrackButton(m.Tracking, m.Parked)
}
