package main

import (
	"log/slog"
	"net/http"
)

func (app *application) websocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("upgrade", "error", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			slog.Error("WS read", "error", err)
			break
		}
		slog.Debug("WS recv", "message", message)

		err = c.WriteMessage(mt, message)
		if err != nil {
			slog.Error("WS write", "error", err)
			break
		}
	}
}
