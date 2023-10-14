package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	indiserver "github.com/adriffaud/indi-web/internal/indi-server"
	"github.com/beevik/etree"
	"github.com/julienschmidt/httprouter"
)

type Device struct {
	Name          string
	Manufacturer  string
	DriverCaption string
	DriverName    string
	Version       string
}

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
	files, err := os.ReadDir("/usr/share/indi")
	if err != nil {
		log.Printf("could not list INDI drivers: %v", err)
		return
	}

	groups := make(map[string][]Device)

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".xml") {
			doc := etree.NewDocument()
			err = doc.ReadFromFile(fmt.Sprintf("/usr/share/indi/%s", file.Name()))
			if err != nil {
				log.Printf("could not list INDI drivers: %v", err)
				return
			}

			driversList := doc.SelectElement("driversList")
			if driversList == nil {
				continue
			}

			for _, groupElem := range driversList.SelectElements("devGroup") {
				var drivers []Device

				group := groupElem.SelectAttrValue("group", "")

				if existingDrivers, ok := groups[group]; ok {
					drivers = append(existingDrivers, drivers...)
				}

				for _, driver := range groupElem.ChildElements() {
					driverChild := driver.SelectElement("driver")

					device := Device{
						Name:          driver.SelectAttrValue("label", ""),
						Manufacturer:  driver.SelectAttrValue("manufacturer", ""),
						DriverCaption: driverChild.SelectAttrValue("name", ""),
						DriverName:    driverChild.Text(),
						Version:       driver.SelectElement("version").Text(),
					}

					drivers = append(drivers, device)
				}

				groups[group] = drivers
			}
		}
	}

	for _, drivers := range groups {
		sort.Slice(drivers, func(i, j int) bool {
			return drivers[i].Name < drivers[j].Name
		})
	}

	jsonb, _ := json.Marshal(groups)
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
