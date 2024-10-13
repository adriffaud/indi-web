package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/adriffaud/indi-web/components"
	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
)

type SSEClient struct {
	app       *application
	eventChan chan string
	request   *http.Request
}

func (sse SSEClient) OnNotify(e indiclient.Event) {
	var buf bytes.Buffer
	var err error

	switch e.EventType {
	case indiclient.Add, indiclient.Delete:
		err = components.DeviceView(sse.app.indiClient.Properties, e.Property.Device).Render(sse.request.Context(), &buf)
	case indiclient.Update:
		err = components.PropertyValues(e.Property).Render(sse.request.Context(), &buf)
	case indiclient.Message:
		slog.Debug("📮 NOTIFICATION", "event", e)
		err = components.Notifications(e.Message).Render(context.Background(), &buf)
	}

	if err != nil {
		slog.Error("failed to convert to HTML", "error", err)
	}

	sse.eventChan <- buf.String()
}

func (app *application) sse(w http.ResponseWriter, r *http.Request) {
	slog.Debug("SSE client connected", "address", r.RemoteAddr)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	client := SSEClient{app: app, request: r, eventChan: make(chan string)}
	app.indiClient.Register(client)

	for {
		select {
		case <-r.Context().Done():
			slog.Debug("SSE client disconnected", "address", r.RemoteAddr)
			app.indiClient.Unregister(client)
			return
		case data := <-client.eventChan:
			fmt.Fprintf(w, "data: %s\n\n", data)
			w.(http.Flusher).Flush()
		}
	}
}
