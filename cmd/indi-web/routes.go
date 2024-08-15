package main

import "net/http"

func (app application) routes() http.Handler {
	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("assets"))
	router.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	router.HandleFunc("GET /", app.index)
	router.HandleFunc("POST /", app.index)
	router.HandleFunc("GET /hardware", app.hardware)
	router.HandleFunc("POST /indi/action", app.indiAction)
	router.HandleFunc("GET /setup", app.setup)
	router.HandleFunc("POST /setup", app.INDIServer)

	router.HandleFunc("GET /ws", app.websocket)
	router.HandleFunc("GET /sse", app.sse)

	return app.checkServerStarted(router)
}
