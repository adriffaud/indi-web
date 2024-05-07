package main

import (
	"log"
	"net/http"
	"time"

	"github.com/adriffaud/indi-web/components"
	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
	"github.com/julienschmidt/httprouter"
)

func (app *application) index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !indiserver.IsRunning() {
		http.Redirect(w, r, "/setup", http.StatusTemporaryRedirect)
	}

	log.Printf("INDI Client: %+v\n", app.indiClient)

	components.Main(indiserver.IsRunning()).Render(r.Context(), w)
}

func (app *application) setup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	driverGroups, err := indiserver.ListDrivers()
	if err != nil {
		log.Printf("could not get INDI drivers: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("INDI Client: %+v\n", app.indiClient)

	devices := make(map[string]indiserver.Device)
	for _, driver := range driverGroups["Telescopes"] {
		if driver.DriverName == "indi_simulator_telescope" && driver.Manufacturer == "Simulator" {
			devices["mount"] = driver
		}
	}
	for _, driver := range driverGroups["CCDs"] {
		if driver.DriverName == "indi_simulator_ccd" && driver.Manufacturer == "Simulator" {
			devices["ccd"] = driver
		} else if driver.DriverName == "indi_simulator_guide" && driver.Manufacturer == "Simulator" {
			devices["guide"] = driver
		}
	}

	log.Printf("%+v\n", devices)

	components.Setup(driverGroups, devices).Render(r.Context(), w)
}

func (app *application) INDIServer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if indiserver.IsRunning() {
		err := indiserver.Stop()
		if err != nil {
			log.Printf("could not stop INDI server: %v", err)
			return
		}

		defer app.indiClient.Close()

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
	time.Sleep(40 * time.Millisecond)

	client, err := indiclient.New("localhost:7624")
	app.indiClient = *client
	if err != nil {
		log.Printf("could not start INDI client: %v", err)
		return
	}

	app.indiClient.GetProperties()

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
