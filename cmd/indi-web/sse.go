package main

import (
	"log/slog"
	"net/http"
)

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
