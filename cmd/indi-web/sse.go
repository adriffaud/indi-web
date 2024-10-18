package main

import (
	"log/slog"
	"net/http"
)

// func (sse SSEClient) OnNotify(e indiclient.Event) {
// 	var component templ.Component
//
// 	switch e.EventType {
// 	case indiclient.Add, indiclient.Delete:
// 		component = components.DeviceView(sse.indiClient.Properties, e.Property.Device)
// 	case indiclient.Update:
// 		component = components.PropertyValues(e.Property)
// 	case indiclient.Message:
// 		slog.Debug("📮 Notification", "message", e.Message)
// 		component = components.Notifications(e.Message)
// 	}
//
// 	if component == nil {
// 		return
// 	}
//
// 	sse.eventChan <- component
// }

func (app *application) sse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")

	app.sseConnectionManager.Register(w, r)

	for {
		select {
		case <-r.Context().Done():
			slog.Debug("SSE client disconnected", "address", r.RemoteAddr)
			app.sseConnectionManager.Unregister(r.RemoteAddr)
			return
		}
	}
}
