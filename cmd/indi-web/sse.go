package main

import (
	"log/slog"
	"net/http"

	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
)

type SSE struct {
}

func (sse SSE) OnNotify(e indiclient.Event) {
	slog.Debug("Received event", "event", e)
}

func (app *application) sse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}
