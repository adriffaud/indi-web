package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
	"github.com/gorilla/websocket"
)

type Websocket struct {
	socket *websocket.Conn
}

func (ws Websocket) OnNotify(e indiclient.Event) {
	js, err := json.Marshal(e)
	if err != nil {
		slog.Error("could not serialize event", "error", err)
	}

	ws.socket.WriteMessage(websocket.TextMessage, js)
}

func (app *application) websocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("upgrade", "error", err)
		return
	}
	defer c.Close()

	ws := Websocket{socket: c}
	app.indiClient.Register(ws)

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
