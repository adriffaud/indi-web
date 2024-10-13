package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/adriffaud/indi-web/components"
	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
)

type SSEClient struct {
	eventChan chan indiclient.Event
}

func (sse SSEClient) OnNotify(e indiclient.Event) {
	sse.eventChan <- e
}

func (app *application) sse(w http.ResponseWriter, r *http.Request) {
	slog.Debug("✅ SSE client connected", "address", r.RemoteAddr, "referer", r.Referer())

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	client := SSEClient{eventChan: make(chan indiclient.Event)}
	app.indiClient.Register(client)

	for {
		select {
		case <-r.Context().Done():
			slog.Debug("SSE client disconnected", "address", r.RemoteAddr)
			app.indiClient.Unregister(client)
			return
		case e := <-client.eventChan:
			var buf bytes.Buffer
			var err error

			switch e.EventType {
			case indiclient.Add, indiclient.Delete:
				err = components.DeviceView(app.indiClient.Properties, e.Property.Device).Render(r.Context(), &buf)
			case indiclient.Update:
				err = components.PropertyValues(e.Property).Render(r.Context(), &buf)

				if e.Property.Device == app.mount.Driver && e.Property.Name == "EQUATORIAL_EOD_COORD" {
					var appBuf bytes.Buffer

					for _, value := range e.Property.Values {
						if value.Name == "RA" {
							app.mount.RA = value.Value
							err := components.TextInput(
								"ra_input",
								app.mount.RA,
								templ.Attributes{"disabled": "true", "hx-swap-oob": "true"},
							).Render(r.Context(), &appBuf)
							if err != nil {
								slog.Error("failed to convert to HTML", "error", err)
							}
						}
						if value.Name == "DEC" {
							app.mount.DEC = value.Value
							err := components.TextInput(
								"dec_input",
								app.mount.DEC,
								templ.Attributes{"disabled": "true", "hx-swap-oob": "true"},
							).Render(r.Context(), &appBuf)
							if err != nil {
								slog.Error("failed to convert to HTML", "error", err)
							}
						}

						fmt.Fprint(w, "event: Mount\n")
						fmt.Fprintf(w, "data: %s\n\n", appBuf.String())
						w.(http.Flusher).Flush()
					}
				}
			case indiclient.Message:
				slog.Debug("📮 Notification", "message", e.Message)
				err = components.Notifications(e.Message).Render(r.Context(), &buf)
			}

			if err != nil {
				slog.Error("failed to convert to HTML", "error", err)
			}
			fmt.Fprint(w, "event: Hardware\n")
			fmt.Fprintf(w, "data: %s\n\n", buf.String())
			w.(http.Flusher).Flush()
		}
	}
}
