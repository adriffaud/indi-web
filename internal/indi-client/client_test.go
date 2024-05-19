package indiclient

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefProperty(t *testing.T) {
	client := &Client{Properties: make([]Property, 0)}

	elements := `
	<defNumberVector device="Telescope Simulator" name="MOUNT_AXES" label="Mount Axes" group="Simulation" state="Idle" perm="ro" timeout="0" timestamp="2024-05-16T12:21:52">
		<defNumber name="PRIMARY" label="Primary (Ha)" format="%g" min="-180" max="180" step="0.010000000000000000208">
		0
		</defNumber>
		<defNumber name="SECONDARY" label="Secondary (Dec)" format="%g" min="-180" max="180" step="0.010000000000000000208">
		0
		</defNumber>
	</defNumberVector>
	<defSwitchVector device="Telescope Simulator" name="SIM_PIER_SIDE" label="Sim Pier Side" group="Simulation" state="Idle" perm="wo" rule="OneOfMany" timeout="60" timestamp="2024-05-16T12:21:52">
	    <defSwitch name="PS_OFF" label="Off">
	Off
	    </defSwitch>
	    <defSwitch name="PS_ON" label="On">
	On
	    </defSwitch>
	</defSwitchVector>

	<defTextVector device="Telescope Simulator" name="ACTIVE_DEVICES" label="Snoop devices" group="Options" state="Idle" perm="rw" timeout="60" timestamp="2024-05-16T12:21:52">
	    <defText name="ACTIVE_GPS" label="GPS">
	GPS Simulator
	    </defText>
	</defTextVector>
	`
	elementsReader := strings.NewReader(elements)
	client.listen(elementsReader)

	numberProp := Property{
		Device:    "Telescope Simulator",
		Group:     "Simulation",
		Type:      Number,
		Name:      "MOUNT_AXES",
		Label:     "Mount Axes",
		State:     "Idle",
		Perm:      "ro",
		Timeout:   0,
		Timestamp: "2024-05-16T12:21:52",
		Rule:      "",
		Values:    []Value{{Name: "PRIMARY", Label: "Primary (Ha)", Value: "0"}, {Name: "SECONDARY", Label: "Secondary (Dec)", Value: "0"}},
	}

	switchProp := Property{
		Device:    "Telescope Simulator",
		Group:     "Simulation",
		Type:      Switch,
		Name:      "SIM_PIER_SIDE",
		Label:     "Sim Pier Side",
		State:     "Idle",
		Perm:      "wo",
		Timeout:   60,
		Timestamp: "2024-05-16T12:21:52",
		Rule:      "OneOfMany",
		Values:    []Value{{Name: "PS_OFF", Label: "Off", Value: "Off"}, {Name: "PS_ON", Label: "On", Value: "On"}},
	}

	textProp := Property{
		Device:    "Telescope Simulator",
		Group:     "Options",
		Type:      Text,
		Name:      "ACTIVE_DEVICES",
		Label:     "Snoop devices",
		State:     "Idle",
		Perm:      "rw",
		Timeout:   60,
		Timestamp: "2024-05-16T12:21:52",
		Rule:      "",
		Values:    []Value{{Name: "ACTIVE_GPS", Label: "GPS", Value: "GPS Simulator"}},
	}

	expected := []Property{numberProp, switchProp, textProp}

	assert.Equal(t, 3, len(client.Properties))
	assert.ElementsMatch(t, expected, client.Properties)
}

func TestUpdateProperty(t *testing.T) {
	numberProp := Property{
		Device:    "Telescope Simulator",
		Group:     "Simulation",
		Type:      Number,
		Name:      "MOUNT_AXES",
		Label:     "Mount Axes",
		State:     "Idle",
		Perm:      "ro",
		Timeout:   0,
		Timestamp: "2024-05-16T12:21:52",
		Rule:      "",
		Values:    []Value{{Name: "PRIMARY", Label: "Primary (Ha)", Value: "0"}, {Name: "SECONDARY", Label: "Secondary (Dec)", Value: "0"}},
	}

	switchProp := Property{
		Device:    "Telescope Simulator",
		Group:     "Simulation",
		Type:      Switch,
		Name:      "SIM_PIER_SIDE",
		Label:     "Sim Pier Side",
		State:     "Idle",
		Perm:      "wo",
		Timeout:   60,
		Timestamp: "2024-05-16T12:21:52",
		Rule:      "OneOfMany",
		Values:    []Value{{Name: "PS_OFF", Label: "Off", Value: "Off"}, {Name: "PS_ON", Label: "On", Value: "On"}},
	}

	textProp := Property{
		Device:    "Telescope Simulator",
		Group:     "Options",
		Type:      Text,
		Name:      "ACTIVE_DEVICES",
		Label:     "Snoop devices",
		State:     "Idle",
		Perm:      "rw",
		Timeout:   60,
		Timestamp: "2024-05-16T12:21:52",
		Rule:      "",
		Values:    []Value{{Name: "ACTIVE_GPS", Label: "GPS", Value: "GPS Simulator"}},
	}

	properties := []Property{numberProp, switchProp, textProp}
	cpy := make([]Property, 3)
	copy(cpy, properties)
	client := &Client{Properties: cpy}

	elements := `
	<defNumberVector device="Telescope Simulator" name="MOUNT_AXES" label="Mount Axes" group="Simulation" state="Idle" perm="ro" timeout="0" timestamp="2025-05-16T12:21:52">
		<defNumber name="PRIMARY" label="Primary (Ha)" format="%g" min="-180" max="180" step="0.010000000000000000208">
		0
		</defNumber>
		<defNumber name="SECONDARY" label="Secondary (Dec)" format="%g" min="-180" max="180" step="0.010000000000000000208">
		0
		</defNumber>
	</defNumberVector>
	<defSwitchVector device="Telescope Simulator" name="SIM_PIER_SIDE" label="Sim Pier Side" group="Simulation" state="Idle" perm="wo" rule="OneOfMany" timeout="60" timestamp="2024-05-16T12:21:52">
	    <defSwitch name="PS_OFF" label="Off">
	Off
	    </defSwitch>
	    <defSwitch name="PS_ON" label="On">
	On
	    </defSwitch>
	</defSwitchVector>

	<defTextVector device="Telescope Simulator" name="ACTIVE_DEVICES" label="Snoop devices" group="Options" state="Idle" perm="rw" timeout="60" timestamp="2024-05-16T12:21:52">
	    <defText name="ACTIVE_GPS" label="GPS">
	GPS Simulator
	    </defText>
	</defTextVector>
	`
	elementsReader := strings.NewReader(elements)
	client.listen(elementsReader)

	numberProp.Timestamp = "2025-05-16T12:21:52"
	expected := []Property{numberProp, switchProp, textProp}
	assert.Equal(t, 3, len(client.Properties))
	assert.ElementsMatch(t, expected, client.Properties)
}

func TestDelProperty(t *testing.T) {
	switchProp := Property{
		Device:    "Telescope Simulator",
		Group:     "Simulation",
		Type:      Switch,
		Name:      "TELESCOPE_PIER_SIDE",
		Label:     "Telescope Pier Side",
		State:     "Idle",
		Perm:      "wo",
		Timeout:   60,
		Timestamp: "2024-05-16T12:21:52",
		Rule:      "OneOfMany",
		Values:    []Value{{Name: "PS_OFF", Label: "Off", Value: "Off"}, {Name: "PS_ON", Label: "On", Value: "On"}},
	}

	properties := []Property{switchProp}
	client := &Client{Properties: properties}

	elements := `
	<delProperty device="Telescope Simulator" name="TELESCOPE_PIER_SIDE" timestamp="2024-05-16T12:47:39"/>
	`
	elementsReader := strings.NewReader(elements)
	client.listen(elementsReader)

	assert.Equal(t, 0, len(client.Properties))
}

func TestSetPropertyValues(t *testing.T) {
	numberProp := Property{
		Device:    "Telescope Simulator",
		Name:      "EQUATORIAL_EOD_COORD",
		State:     "Idle",
		Timestamp: "2024-05-16T12:40:10",
		Values:    []Value{{Name: "RA", Value: "0"}, {Name: "DEC", Value: "0"}},
	}

	properties := []Property{numberProp}
	client := &Client{Properties: properties}

	elements := `
	<setNumberVector device="Telescope Simulator" name="EQUATORIAL_EOD_COORD" state="Idle" timeout="60" timestamp="2024-05-16T12:48:10">
		<oneNumber name="RA">
			22.451127260193981527
		</oneNumber>
		<oneNumber name="DEC">
			90
		</oneNumber>
	</setNumberVector>
	`
	elementsReader := strings.NewReader(elements)
	client.listen(elementsReader)

	expected := Property{
		Device:    "Telescope Simulator",
		Name:      "EQUATORIAL_EOD_COORD",
		State:     "Idle",
		Timestamp: "2024-05-16T12:48:10",
		Values:    []Value{{Name: "RA", Value: "22.451127260193981527"}, {Name: "DEC", Value: "90"}},
	}
	assert.Equal(t, 1, len(client.Properties))
	assert.Equal(t, expected, client.Properties[0])
}

func TestSetUnexistingPropertyValues(t *testing.T) {
	client := &Client{Properties: make(Properties, 0)}

	elements := `
	<setNumberVector device="Telescope Simulator" name="EQUATORIAL_EOD_COORD" state="Idle" timeout="60" timestamp="2024-05-16T12:48:10">
		<oneNumber name="RA">
			22.451127260193981527
		</oneNumber>
		<oneNumber name="DEC">
			90
		</oneNumber>
	</setNumberVector>
	`
	elementsReader := strings.NewReader(elements)

	assert.Panics(t, func() { client.listen(elementsReader) })
}
