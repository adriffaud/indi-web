package indiclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDevicesSorted(t *testing.T) {
	properties := Properties{Property{Device: "device1"}, Property{Device: "device1"}, Property{Device: "device2"}}
	devices := properties.GetDevicesSorted()
	expected := []string{"device1", "device2"}

	assert.Equal(t, 2, len(devices))
	assert.ElementsMatch(t, expected, devices)
}

func TestGetDevicesGroupsSorted(t *testing.T) {
	properties := Properties{Property{Device: "device1", Group: "group2"}, Property{Device: "device1", Group: "group1"}, Property{Device: "device2", Group: "unwanted"}}
	groups := properties.GetDeviceGroupsSorted("device1")
	expected := []string{"group1", "group2"}

	assert.Equal(t, 2, len(groups))
	assert.ElementsMatch(t, expected, groups)
}

func TestGetPropertiesForDeviceGroup(t *testing.T) {
	properties := Properties{Property{Device: "device1", Group: "group2"}, Property{Device: "device1", Group: "group1", Name: "foo"}, Property{Device: "device1", Group: "group1", Name: "bar"}, Property{Device: "device2", Group: "unwanted"}}
	filteredProperties := properties.GetPropertiesForDeviceGroup("device1", "group1")
	expected := []Property{{Device: "device1", Group: "group1", Name: "foo"}, {Device: "device1", Group: "group1", Name: "bar"}}

	assert.Equal(t, 2, len(filteredProperties))
	assert.ElementsMatch(t, expected, filteredProperties)
}
