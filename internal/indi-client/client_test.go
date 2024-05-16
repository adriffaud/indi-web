package indiclient

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefProperty(t *testing.T) {
	client := &Client{Properties: make([]Property, 0)}

	defNumberVector := `
	<defNumberVector device="Telescope Simulator" name="MOUNT_AXES" label="Mount Axes" group="Simulation" state="Idle" perm="ro" timeout="0" timestamp="2024-05-16T12:21:52">
		<defNumber name="PRIMARY" label="Primary (Ha)" format="%g" min="-180" max="180" step="0.010000000000000000208">
		0
		</defNumber>
		<defNumber name="SECONDARY" label="Secondary (Dec)" format="%g" min="-180" max="180" step="0.010000000000000000208">
		0
		</defNumber>
	</defNumberVector>`
	defNumberReader := strings.NewReader(defNumberVector)
	client.listen(defNumberReader)

	expected := Property{
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

	assert.Equal(t, 1, len(client.Properties))
	assert.Equal(t, expected, client.Properties[0])
}

// <defSwitchVector device="Telescope Simulator" name="SIM_PIER_SIDE" label="Sim Pier Side" group="Simulation" state="Idle" perm="wo" rule="OneOfMany" timeout="60" timestamp="2024-05-16T12:21:52">
//     <defSwitch name="PS_OFF" label="Off">
// Off
//     </defSwitch>
//     <defSwitch name="PS_ON" label="On">
// On
//     </defSwitch>
// </defSwitchVector>

// <defTextVector device="Telescope Simulator" name="ACTIVE_DEVICES" label="Snoop devices" group="Options" state="Idle" perm="rw" timeout="60" timestamp="2024-05-16T12:21:52">
//     <defText name="ACTIVE_GPS" label="GPS">
// GPS Simulator
//     </defText>
//     <defText name="ACTIVE_DOME" label="DOME">
// Dome Simulator
//     </defText>
// </defTextVector>

// <delProperty device="Telescope Simulator" name="TELESCOPE_PIER_SIDE" timestamp="2024-05-16T12:47:39"/>

// <setNumberVector device="Telescope Simulator" name="EQUATORIAL_EOD_COORD" state="Idle" timeout="60" timestamp="2024-05-16T12:48:10">
//     <oneNumber name="RA">
// 22.451127260193981527
//     </oneNumber>
//     <oneNumber name="DEC">
// 90
//     </oneNumber>
// </setNumberVector>
