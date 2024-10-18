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
	}

	if event.EventType != indiclient.Update || m.Driver == "" || event.Property.Device != m.Driver {
		return
	}

	if event.Property.Name == "EQUATORIAL_EOD_COORD" {
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
	}

}

func (m *Mount) SetClient(client *indiclient.Client) {
	m.client = client
}

func (m Mount) Park() {}

func (m Mount) Unpark() {}

func (m Mount) StartTracking() {}

func (m Mount) StopTracking() {}
