package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adriffaud/indi-web/components"
	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	driverGroups, err := indiserver.ListDrivers()
	if err != nil {
		log.Printf("could not get INDI drivers: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	components.Page(indiserver.IsRunning(), driverGroups).Render(r.Context(), w)
}

func INDIDrivers(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	driverGroups, err := indiserver.ListDrivers()
	if err != nil {
		log.Printf("could not get INDI drivers: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonb, _ := json.Marshal(driverGroups)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonb))
}

func INDIServer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if indiserver.IsRunning() {
		err := indiserver.Stop()
		if err != nil {
			log.Printf("could not stop INDI server: %v", err)
			return
		}
		fmt.Fprint(w, "Start")
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	log.Printf("%+v\n", r.Form)
	d := make([]string, 0, len(r.Form))
	for _, driver := range r.Form {
		d = append(d, driver[0])
	}

	err = indiserver.Start(d)
	if err != nil {
		log.Printf("could not start INDI server: %v", err)
		return
	}
	fmt.Fprint(w, "Stop")
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/indi/activate", INDIServer)
	router.GET("/indi/drivers", INDIDrivers)
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
