package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
)

type application struct {
	indiClient *indiclient.Client
}

var (
	host string
	port int
)

func main() {
	flag.StringVar(&host, "host", "localhost", "server host")
	flag.IntVar(&port, "port", 8080, "server port")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	app := &application{}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      app.routes(),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	slog.Info("starting server", "addr", server.Addr)
	server.ListenAndServe()
}
