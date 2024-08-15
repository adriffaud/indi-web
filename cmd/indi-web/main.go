package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
	"github.com/gorilla/websocket"
)

type application struct {
	indiClient *indiclient.Client
}

var (
	host     string
	port     int
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	flag.StringVar(&host, "host", "", "server host")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	app := &application{}

	// TODO: REMOVE ME
	// TEMP AUTOSTART
	err := indiserver.Start([]string{"indi_simulator_telescope"})
	if err != nil {
		slog.Info("could not start INDI server", "error", err)
		return
	}

	time.Sleep(40 * time.Millisecond)

	client, err := indiclient.New("localhost:7624")
	if err != nil {
		slog.Info("could not start INDI client", "error", err)
		return
	}
	app.indiClient = client
	app.indiClient.GetProperties()

	slog.Debug("INDI client connected")
	// TEMP AUTOSTART

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: app.routes(),
	}

	slog.Info("starting server", "addr", server.Addr)
	server.ListenAndServe()
}
