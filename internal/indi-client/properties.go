package indiclient

import (
	"sort"
)

type PropertyType int64

const (
	Number PropertyType = iota
	Switch
	Text
)

type Value struct {
	Name  string
	Label string
	Value string
}

type Property struct {
	Device    string
	Group     string
	Type      PropertyType
	Name      string
	Label     string
	State     string
	Perm      string
	Timeout   int
	Timestamp string
	Rule      string
	Values    []Value
}

type Properties []Property

func (p Properties) GetDevicesSorted() []string {
	allKeys := make(map[string]bool)
	devices := make([]string, 0)
	for _, property := range p {
		if _, value := allKeys[property.Device]; !value {
			allKeys[property.Device] = true
			devices = append(devices, property.Device)
		}
	}

	sort.Strings(devices)
	return devices
}

func (p Properties) GetDeviceGroupsSorted(device string) []string {
	filtered := make(Properties, 0)
	for _, property := range p {
		if property.Device == device {
			filtered = append(filtered, property)
		}
	}

	allKeys := make(map[string]bool)
	groups := make([]string, 0)
	for _, property := range filtered {
		if _, value := allKeys[property.Group]; !value {
			allKeys[property.Group] = true
			groups = append(groups, property.Group)
		}
	}

	return groups
}

func (p Properties) GetPropertiesForDeviceGroup(device, group string) Properties {
	filtered := make(Properties, 0)
	for _, property := range p {
		if property.Device == device && property.Group == group {
			filtered = append(filtered, property)
		}
	}
	return filtered
}

func (p Properties) FindProperty(selector PropertySelector) *Property {
	for _, property := range p {
		if property.Device == selector.Device && property.Name == selector.Name {
			return &property
		}
	}

	return nil
}
