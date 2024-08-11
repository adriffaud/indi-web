package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app application) routes() http.Handler {
	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("assets"))

	router.GET("/", app.index)
	router.POST("/", app.index)
	router.GET("/hardware", app.hardware)
	router.POST("/indi/action", app.indiAction)
	router.GET("/setup", app.setup)
	router.POST("/setup", app.INDIServer)

	return app.checkServerStarted(router)
}
