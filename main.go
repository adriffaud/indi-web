package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adriffaud/indi-web/components"
	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
	indiclient "github.com/adriffaud/indi-web/pkg/indi-client"
	"github.com/julienschmidt/httprouter"
)

var indiClient *indiclient.Client

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	driverGroups, err := indiserver.ListDrivers()
	if err != nil {
		log.Printf("could not get INDI drivers: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("INDI Client: %+v\n", indiClient)

	components.Page(indiserver.IsRunning(), driverGroups).Render(r.Context(), w)
}

func INDIServer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if indiserver.IsRunning() {
		err := indiserver.Stop()
		if err != nil {
			log.Printf("could not stop INDI server: %v", err)
			return
		}
		w.Header().Add("HX-Location", "/")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	d := make([]string, 0, len(r.Form))
	for _, driver := range r.Form {
		d = append(d, driver[0])
	}

	err = indiserver.Start(d)
	if err != nil {
		log.Printf("could not start INDI server: %v", err)
		return
	}

	// TODO: Wait for server start before creating the client
	time.Sleep(400 * time.Millisecond)

	indiClient, err = indiclient.New("localhost:7624")
	if err != nil {
		log.Printf("could not start INDI client: %v", err)
		return
	}

	components.IndiServerButton(true).Render(r.Context(), w)
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/indi/activate", INDIServer)
	router.ServeFiles("/static/*filepath", http.Dir("assets"))

	server := &http.Server{
		Addr:         "localhost:8080",
		Handler:      router,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	fmt.Printf("Listening on %v\n", server.Addr)
	server.ListenAndServe()
}
