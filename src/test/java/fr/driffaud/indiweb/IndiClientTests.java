package fr.driffaud.indiweb;

import org.junit.jupiter.api.Test;

import java.io.BufferedReader;
import java.io.StringReader;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;

public class IndiClientTests {

    @Test
    void TestDefOneProperty() {
        var element = """
                <defNumberVector device="Telescope Simulator" name="MOUNT_AXES" label="Mount Axes" group="Simulation" state="Idle" perm="ro" timeout="0" timestamp="2024-05-16T12:21:52">
                	<defNumber name="PRIMARY" label="Primary (Ha)" format="%g" min="-180" max="180" step="0.010000000000000000208">
                	0
                	</defNumber>
                	<defNumber name="SECONDARY" label="Secondary (Dec)" format="%g" min="-180" max="180" step="0.010000000000000000208">
                	0
                	</defNumber>
                </defNumberVector>
                """;
        var inputString = new StringReader(element);
        var reader = new BufferedReader(inputString);

        var client = new IndiClient();
        client.listen(reader);

        var expected = List.of(new Property("Telescope Simulator", "Simulation",
                PropertyType.NUMBER, "MOUNT_AXES", "Mount Axes", "Idle", "ro",
                List.of(new Value("PRIMARY", "Primary (Ha)", "0"), new Value("SECONDARY", "Secondary (Dec)", "0"))));

        System.out.println(client.properties);
        assertEquals(1, client.properties.size());
        assertEquals(expected, client.properties);
    }

    @Test
    void TestDefProperties() {
        var elements = """
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
                """;
        var inputString = new StringReader(elements);
        var reader = new BufferedReader(inputString);

        var client = new IndiClient();
        client.listen(reader);

        var expected = List.of(new Property("Telescope Simulator", "Simulation",
                        PropertyType.NUMBER, "MOUNT_AXES", "Mount Axes", "Idle", "ro",
                        List.of(new Value("PRIMARY", "Primary (Ha)", "0"), new Value("SECONDARY", "Secondary (Dec)", "0"))),
                new Property("Telescope Simulator", "Simulation",
                        PropertyType.SWITCH, "SIM_PIER_SIDE", "Sim Pier Side", "Idle", "wo",
                        List.of(new Value("PS_OFF", "Off", "Off"), new Value("PS_ON", "On", "On"))),
                new Property("Telescope Simulator", "Options",
                        PropertyType.TEXT, "ACTIVE_DEVICES", "Snoop devices", "Idle", "rw",
                        List.of(new Value("ACTIVE_GPS", "GPS", "GPS Simulator"))));

        System.out.println(client.properties);
        assertEquals(3, client.properties.size());
        assertEquals(expected, client.properties);
    }
}
