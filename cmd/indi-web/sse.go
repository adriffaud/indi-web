package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/adriffaud/indi-web/components"
	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
)

type SSEClient struct {
	addr      string
	eventChan chan indiclient.Event
}

func (sse SSEClient) OnNotify(e indiclient.Event) {
	sse.eventChan <- e
}

func (app *application) sse(w http.ResponseWriter, r *http.Request) {
	slog.Debug("SSE client connected", "address", r.RemoteAddr)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	client := SSEClient{addr: r.RemoteAddr, eventChan: make(chan indiclient.Event)}
	app.indiClient.Register(client)

	for {
		select {
		case <-r.Context().Done():
			slog.Debug("SSE client disconnected", "address", r.RemoteAddr)
			app.indiClient.Unregister(client)
			return
		case evt := <-client.eventChan:
			switch evt.EventType {
			case indiclient.Add, indiclient.Delete:
				tmpl, err := templ.ToGoHTML(r.Context(), components.DeviceView(app.indiClient.Properties, evt.Property.Device))
				if err != nil {
					slog.Error("failed to convert to HTML", "error", err)
				}

				fmt.Fprintf(w, "data: %s\n\n", tmpl)
			case indiclient.Update:
				tmpl, err := templ.ToGoHTML(r.Context(), components.PropertyValues(evt.Property))
				if err != nil {
					slog.Error("failed to convert to HTML", "error", err)
				}

				fmt.Fprintf(w, "data: %s\n\n", tmpl)
			case indiclient.Message:
				slog.Debug("NOTIFICATION", "event", evt)
				var msg bytes.Buffer
				components.Notifications(evt.Message).Render(context.Background(), &msg)
				slog.Debug("HTML", "msg", msg.String())
				fmt.Fprintf(w, "data: %s\n\n", msg.String())
			}

			w.(http.Flusher).Flush()
		}
	}
}
