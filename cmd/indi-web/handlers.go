package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/adriffaud/indi-web/components"
	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
	"github.com/julienschmidt/httprouter"
)

func (app *application) index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	components.Main().Render(r.Context(), w)
}

func (app *application) hardware(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	components.Hardware(app.indiClient.Properties).Render(r.Context(), w)
}

func (app *application) setup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	driverGroups, err := indiserver.ListDrivers()
	if err != nil {
		slog.Error("could not get INDI drivers", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	devices := make(map[string]indiserver.Device)
	for _, driver := range driverGroups["Telescopes"] {
		if driver.DriverName == "indi_simulator_telescope" && driver.Manufacturer == "Simulator" {
			devices["mount"] = driver
		}
	}
	// for _, driver := range driverGroups["CCDs"] {
	// 	if driver.DriverName == "indi_simulator_ccd" && driver.Manufacturer == "Simulator" {
	// 		devices["ccd"] = driver
	// 	} else if driver.DriverName == "indi_simulator_guide" && driver.Manufacturer == "Simulator" {
	// 		devices["guide"] = driver
	// 	}
	// }
	// for _, driver := range driverGroups["Focusers"] {
	// 	if driver.DriverName == "indi_simulator_focus" && driver.Manufacturer == "Simulator" {
	// 		devices["focuser"] = driver
	// 	}
	// }

	components.Setup(driverGroups, devices).Render(r.Context(), w)
}

func (app *application) INDIServer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	slog.Debug("Handling INDI server setup")

	if indiserver.IsRunning() {
		slog.Debug("Stopping INDI server")
		err := indiserver.Stop()
		if err != nil {
			slog.Error("could not stop INDI server", "error", err)
			return
		}

		defer app.indiClient.Close()

		http.Redirect(w, r, "/setup", http.StatusSeeOther)
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
		slog.Info("could not start INDI server", "error", err)
		return
	}

	// TODO: Wait for server start before creating the client
	time.Sleep(40 * time.Millisecond)

	client, err := indiclient.New("localhost:7624")
	app.indiClient = client
	if err != nil {
		slog.Info("could not start INDI client", "error", err)
		return
	}
	app.indiClient.GetProperties()

	slog.Debug("INDI client connected")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
