package indiserver

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/beevik/etree"
)

// Device represents a device in the drivers list.
type Device struct {
	Name          string
	Manufacturer  string
	DriverCaption string
	DriverName    string
	Version       string
}

type DeviceGroups map[string][]Device

// ListDrivers returns a map of grouped devices.
func ListDrivers() (DeviceGroups, error) {
	files, err := os.ReadDir("/usr/share/indi")
	if err != nil {
		return nil, err
	}

	groups := make(DeviceGroups)

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".xml") {
			err = parseXML(fmt.Sprintf("/usr/share/indi/%s", file.Name()), groups)
			if err != nil {
				return nil, err
			}
		}
	}

	sortDrivers(groups)

	return groups, nil
}

func parseXML(filepath string, groups DeviceGroups) error {
	doc := etree.NewDocument()
	err := doc.ReadFromFile(filepath)
	if err != nil {
		return err
	}

	driversList := doc.SelectElement("driversList")
	if driversList == nil {
		return nil
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

	return nil
}

func sortDrivers(groups DeviceGroups) {
	for _, drivers := range groups {
		sort.Slice(drivers, func(i, j int) bool {
			return drivers[i].Name < drivers[j].Name
		})
	}
}
