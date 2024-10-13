package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/adriffaud/indi-web/internal/config"
	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
)

type application struct {
	indiClient *indiclient.Client
	mount      config.Mount
}

func (app application) OnNotify(e indiclient.Event) {
	switch e.EventType {
	case indiclient.Timeout:
		if app.mount.Driver != "" && !app.mount.Connected {
			slog.Debug("🤓 Mount not connected, connecting...")
			err := app.indiClient.Connect(app.mount.Driver)
			if err != nil {
				slog.Error("🔴 Could not automatically connect to mount", "err", err)
			}
		}
	case indiclient.Update:
		if e.Property.Device == app.mount.Driver && e.Property.Name == "EQUATORIAL_EOD_COORD" {
			for _, value := range e.Property.Values {
				if value.Name == "RA" {
					app.mount.RA = value.Value
				}
				if value.Name == "DEC" {
					app.mount.DEC = value.Value
				}
			}
		}
	}
}

var (
	host string
	port int
)

func main() {
	flag.StringVar(&host, "host", "", "server host")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	app := &application{}
	app.mount.Driver = "Telescope Simulator"

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
	app.indiClient.Register(app)

	slog.Debug("INDI client connected")
	// TEMP AUTOSTART

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: app.routes(),
	}

	slog.Info("starting server", "addr", server.Addr)
	server.ListenAndServe()
}
