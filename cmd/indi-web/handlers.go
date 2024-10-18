package main

import (
	"log/slog"
	"net/http"
	"time"

	indiclient "github.com/adriffaud/indi-web/internal/indi-client"
	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
	"github.com/adriffaud/indi-web/ui/pages"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	pages.Main(app.mount).Render(r.Context(), w)
}

func (app *application) hardware(w http.ResponseWriter, r *http.Request) {
	pages.Hardware(app.indiClient.Properties).Render(r.Context(), w)
}

func (app *application) setup(w http.ResponseWriter, r *http.Request) {
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

	pages.Setup(driverGroups, devices).Render(r.Context(), w)
}

func (app *application) indiAction(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "unable to parse form", http.StatusBadRequest)
		return
	}

	// TODO: Add form validation and struct mapping
	selector := indiclient.PropertySelector{
		Device:    r.FormValue("device"),
		Name:      r.FormValue("name"),
		ValueName: r.FormValue("valueName"),
	}

	slog.Debug("Sending connection property", "selector", selector)
	err = app.indiClient.NewPropertyValue(selector)
	if err != nil {
		slog.Error("INDI action", "selector", selector, "error", err)
		http.Error(w, "could not send new property value", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *application) INDIServer(w http.ResponseWriter, r *http.Request) {
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

	client, err := indiclient.New("localhost:7624", app.eventChan)
	app.indiClient = client
	if err != nil {
		slog.Info("could not start INDI client", "error", err)
		return
	}
	app.indiClient.GetProperties()

	slog.Debug("INDI client connected")

	http.Redirect(w, r, "/hardware", http.StatusSeeOther)
}
