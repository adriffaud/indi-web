package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	tmpl, _ := os.ReadFile("web/template/index.html")

	driverGroups, err := indiserver.ListDrivers()
	if err != nil {
		log.Printf("could not get INDI drivers: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		DriversGroups indiserver.DeviceGroups
		Running       bool
	}{
		DriversGroups: driverGroups,
		Running:       indiserver.IsRunning(),
	}
	t, err := template.New("index").Funcs(template.FuncMap{"lower": strings.ToLower, "appendSuffix": func(suffix, str string) string {
		return str + suffix
	}}).Parse(string(tmpl))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	log.Printf("%+v\n", r.Form)

	if indiserver.IsRunning() {
		err := indiserver.Stop()
		if err != nil {
			log.Printf("could not stop INDI server: %v", err)
			return
		}
		fmt.Fprint(w, "Start")
	} else {
		err := indiserver.Start()
		if err != nil {
			log.Printf("could not start INDI server: %v", err)
			return
		}
		fmt.Fprint(w, "Stop")
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/indi/activate", INDIServer)
	router.GET("/indi/drivers", INDIDrivers)
	router.ServeFiles("/static/*filepath", http.Dir("web/static"))

	log.Fatal(http.ListenAndServe(":8080", router))
}
