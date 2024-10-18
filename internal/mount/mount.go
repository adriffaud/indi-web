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
			mount.OnNotify(event)
		}
	}()

	return &mount
}

func (m *Mount) OnNotify(event indiclient.Event) {
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
			if value.Name == "RA" {
				formated, err := coordconv.DecimalToSexagesimal(value.Value)
				if err != nil {
					slog.Error("could not convert RA coords in sexagesimal", "err", err, "value", value)
				}
				m.RA = formated
				component := components.TextInput("ra_input", m.RA, templ.Attributes{"disabled": "true", "hx-swap-oob": "true"})
				m.htmlChan <- component
				if err != nil {
					slog.Error("failed to convert to HTML", "error", err)
				}
			}
			if value.Name == "DEC" {
				formated, err := coordconv.DecimalToSexagesimal(value.Value)
				if err != nil {
					slog.Error("could not convert DEC coords in sexagesimal", "err", err, "value", value)
				}
				m.DEC = formated
				component := components.TextInput("dec_input", m.DEC, templ.Attributes{"disabled": "true", "hx-swap-oob": "true"})
				m.htmlChan <- component
				if err != nil {
					slog.Error("failed to convert to HTML", "error", err)
				}
			}
		}
	}

}

func (m *Mount) SetClient(client *indiclient.Client) {
	m.client = client
}

func (m Mount) Connect() {}

func (m Mount) Disconnect() {}

func (m Mount) Park() {}

func (m Mount) Unpark() {}

func (m Mount) StartTracking() {}

func (m Mount) StopTracking() {}
