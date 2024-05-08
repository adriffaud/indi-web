package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
)

type application struct {
	indiClient indiclient.Client
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := application{}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      app.routes(),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	slog.Info("starting server", "addr", server.Addr)
	server.ListenAndServe()
}
