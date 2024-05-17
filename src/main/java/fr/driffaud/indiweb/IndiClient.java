package fr.driffaud.indiweb;

import javax.xml.namespace.QName;
import javax.xml.stream.XMLInputFactory;
import javax.xml.stream.XMLStreamException;
import javax.xml.stream.events.StartElement;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.PrintWriter;
import java.net.Socket;
import java.util.ArrayList;
import java.util.List;

public class IndiClient {
    private Socket socket;
    private PrintWriter out;
    private BufferedReader in;
    public List<Property> properties = new ArrayList<>();

    public void start(String address, int port) throws IOException {
        socket = new Socket(address, port);
        out = new PrintWriter(socket.getOutputStream(), true);
        in = new BufferedReader(new InputStreamReader(socket.getInputStream()));

        Thread.ofVirtual().start(() -> listen(in));
    }

    protected void listen(BufferedReader in) {
        var xmlInputFactory = XMLInputFactory.newInstance();
        try {
            var reader = xmlInputFactory.createXMLEventReader(in);

            while (true) {
                var event = reader.nextEvent();

                if (event.isStartElement()) {
                    StartElement startElement = event.asStartElement();
                    System.out.println("Element: " + startElement.getName().getLocalPart());

                    var device = startElement.getAttributeByName(new QName("device"));
                    var group = startElement.getAttributeByName(new QName("group"));
                    var name = startElement.getAttributeByName(new QName("name"));
                    var label = startElement.getAttributeByName(new QName("label"));
                    var state = startElement.getAttributeByName(new QName("state"));
                    var perm = startElement.getAttributeByName(new QName("perm"));
                    var type = switch (startElement.getName().getLocalPart()) {
                        case "defNumberVector" -> PropertyType.NUMBER;
                        case "defSwitchVector" -> PropertyType.SWITCH;
                        case "defTextVector" -> PropertyType.TEXT;
                        default -> PropertyType.UNKNOWN;
                    };

                    var property = new Property(device.getValue(), group.getValue(), type,
                            name.getValue(), label.getValue(), state.getValue(), perm.getValue());
                    properties.add(property);
                }
            }
        } catch (XMLStreamException e) {
            throw new RuntimeException(e);
        }
    }

    public void getProperties() {
        this.sendMessage("<getProperties version=\"1.7\" />");
    }

    private void sendMessage(String msg) {
        System.out.println("Sending: " + msg);
        out.println(msg);
    }

    public void stopConnection() throws IOException {
        in.close();
        socket.close();
    }
}
