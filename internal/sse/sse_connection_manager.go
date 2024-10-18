package sse

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
)

type SSEClient struct {
	writer  http.ResponseWriter
	request *http.Request
}

type SSEConnectionManager struct {
	clients map[string]SSEClient
}

func NewSSEConnectionManager(htmlChan chan templ.Component) SSEConnectionManager {
	mgr := SSEConnectionManager{clients: make(map[string]SSEClient)}

	go func() {
		for {
			component := <-htmlChan
			mgr.sendHTML(component)
		}
	}()

	return mgr
}

func (mgr SSEConnectionManager) Register(writer http.ResponseWriter, request *http.Request) {
	slog.Debug("✅ Registering SSE client", "address", request.RemoteAddr)
	mgr.clients[request.RemoteAddr] = SSEClient{writer: writer, request: request}
}

func (mgr SSEConnectionManager) Unregister(addr string) {
	slog.Debug("🚮 Unregistering SSE client", "address", addr)
	delete(mgr.clients, addr)
}

func (mgr SSEConnectionManager) sendHTML(component templ.Component) {
	for _, client := range mgr.clients {
		var buf bytes.Buffer
		var err error

		component.Render(client.request.Context(), &buf)
		if err != nil {
			slog.Error("failed to convert to HTML", "error", err)
		}

		fmt.Fprintf(client.writer, "data: %s\n\n", buf.String())
		client.writer.(http.Flusher).Flush()
	}
}
