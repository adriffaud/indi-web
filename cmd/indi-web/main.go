package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	tmpl, _ := os.ReadFile("web/template/index.html")
	data := struct {
		Running bool
	}{
		Running: indiserver.IsRunning(),
	}
	t, err := template.New("index").Parse(string(tmpl))
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

func INDIServer(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	if indiserver.IsRunning() {
		err := indiserver.Stop()
		if err != nil {
			log.Printf("could not stop INDI server: %v", err)
			return
		}
		fmt.Fprint(w, "Stopped")
	} else {
		err := indiserver.Start()
		if err != nil {
			log.Printf("could not start INDI server: %v", err)
			return
		}
		fmt.Fprint(w, "Running")
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
