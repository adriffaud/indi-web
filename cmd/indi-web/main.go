package main

import (
	"fmt"
	"net/http"
	"time"

	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
)

type application struct {
	indiClient indiclient.Client
}

func main() {
	app := application{}

	server := &http.Server{
		Addr:         ":8080",
		Handler:      app.routes(),
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	fmt.Printf("Listening on http://localhost%v\n", server.Addr)
	server.ListenAndServe()
}
