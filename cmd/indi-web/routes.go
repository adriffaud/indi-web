package main

import (
	"net/http"

	"github.com/adriffaud/indi-web/assets"
)

func (app application) routes() http.Handler {
	router := http.NewServeMux()

	fileServer := http.FileServer(http.FS(assets.Content))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	router.HandleFunc("GET /", app.index)
	router.HandleFunc("POST /", app.index)
	router.HandleFunc("GET /hardware", app.hardware)
	router.HandleFunc("POST /indi/action", app.indiAction)
	router.HandleFunc("POST /mount/action", app.mountAction)
	router.HandleFunc("GET /setup", app.setup)
	router.HandleFunc("POST /setup", app.INDIServer)

	router.HandleFunc("GET /sse", app.sse)

	return app.checkServerStarted(router)
}
