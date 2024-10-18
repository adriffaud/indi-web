package debug

import (
	"log/slog"

	"github.com/a-h/templ"
	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
	"github.com/adriffaud/indi-web/ui/pages"
)

// DebugWatcher sends SSE events for every client event to update hardware debug page
func DebugWatcher(client *indiclient.Client, eventChan chan indiclient.Event, htmlChan chan templ.Component) {
	go func() {
		for {
			event := <-eventChan
			var component templ.Component

			switch event.EventType {
			case indiclient.Add, indiclient.Delete:
				component = pages.DeviceView(client.Properties, event.Property.Device)
			case indiclient.Update:
				component = pages.PropertyValues(event.Property)
			case indiclient.Message:
				slog.Debug("📮 Notification", "message", event.Message)
				component = pages.Notifications(event.Message)
			}

			if component == nil {
				return
			}

			htmlChan <- component
		}
	}()
}
