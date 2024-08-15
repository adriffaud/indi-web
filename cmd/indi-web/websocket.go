package main

import (
	"fmt"
	"log/slog"
	"net/http"

	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
	"github.com/gorilla/websocket"
)

type WSClient struct {
	name string
	conn *websocket.Conn
}

func (ws WSClient) OnNotify(e indiclient.Event) {
	message := fmt.Sprintf("<div id=\"%s\">%s</div>", e.Property.Device, e.Property.Device)
	slog.Debug("sending message", "message", message)
	ws.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (app *application) websocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("upgrade", "error", err)
		return
	}
	defer c.Close()

	ws := WSClient{name: r.RemoteAddr, conn: c}
	app.indiClient.Register(ws)
	defer app.indiClient.Unregister(ws)
	slog.Debug("new WS client", "client", ws)

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
