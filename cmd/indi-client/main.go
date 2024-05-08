package main

import (
	"flag"
	"log/slog"
	"os"

	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
)

func main() {
	var host string

	flag.StringVar(&host, "host", "localhost:7624", "INDI server address")
	flag.Parse()

	slog.Info("Connecting to INDI server", "host", host)

	client, err := indiclient.New(host)
	if err != nil {
		slog.Error("could not create INDI client", "error", err)
		os.Exit(1)
	}

	slog.Info("connected")

	err = client.GetProperties()
	if err != nil {
		slog.Error("could not get INDI properties", "error", err)
		os.Exit(1)
	}

	// Wait forever until user kills the process
	select {}
}
