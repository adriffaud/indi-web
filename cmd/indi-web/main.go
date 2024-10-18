package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
	"github.com/adriffaud/indi-web/internal/mount"
	"github.com/adriffaud/indi-web/internal/sse"
)

type application struct {
	indiClient           *indiclient.Client
	sseConnectionManager sse.SSEConnectionManager
	eventChan            chan indiclient.Event
	htmlChan             chan templ.Component
	mount                *mount.Mount
}

var (
	host string
	port int
)

func main() {
	flag.StringVar(&host, "host", "", "server host")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	htmlChan := make(chan templ.Component)
	eventChan := make(chan indiclient.Event)

	app := &application{
		htmlChan:             htmlChan,
		eventChan:            eventChan,
		sseConnectionManager: sse.NewSSEConnectionManager(htmlChan),
		mount:                mount.NewMount("Telescope Simulator", eventChan, htmlChan),
	}

	// TODO: REMOVE ME
	// TEMP AUTOSTART
	err := indiserver.Start([]string{"indi_simulator_telescope"})
	if err != nil {
		slog.Info("could not start INDI server", "error", err)
		return
	}

	time.Sleep(40 * time.Millisecond)

	client, err := indiclient.New("localhost:7624", app.eventChan)
	if err != nil {
		slog.Info("could not start INDI client", "error", err)
		return
	}
	app.indiClient = client
	app.mount.SetClient(app.indiClient)
	app.indiClient.GetProperties()

	// TEMP AUTOSTART

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: app.routes(),
	}

	slog.Info("starting server", "addr", server.Addr)
	server.ListenAndServe()
}
