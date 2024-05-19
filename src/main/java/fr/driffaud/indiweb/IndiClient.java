package fr.driffaud.indiweb;

import com.ctc.wstx.api.WstxInputProperties;
import com.ctc.wstx.stax.WstxInputFactory;

import javax.xml.namespace.QName;
import javax.xml.stream.XMLStreamException;
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
        var xmlInputFactory = WstxInputFactory.newInstance();
        // Crucial to allow multiple root elements
        xmlInputFactory.setProperty(WstxInputProperties.P_INPUT_PARSING_MODE, WstxInputProperties.PARSING_MODE_FRAGMENT);

        try {
            var reader = xmlInputFactory.createXMLEventReader(in);
            Property property = null;

            while (reader.hasNext()) {
                var event = reader.nextEvent();

                if (event.isStartElement()) {
                    var element = event.asStartElement();
                    switch (element.getName().getLocalPart()) {
                        case "defNumberVector":
                        case "defSwitchVector":
                        case "defTextVector":
                            property = new Property();
                            property.device = element.getAttributeByName(new QName("device")).getValue();
                            property.group = element.getAttributeByName(new QName("group")).getValue();
                            property.name = element.getAttributeByName(new QName("name")).getValue();
                            property.label = element.getAttributeByName(new QName("label")).getValue();
                            property.state = element.getAttributeByName(new QName("state")).getValue();
                            property.perm = element.getAttributeByName(new QName("perm")).getValue();
                            property.timeout = element.getAttributeByName(new QName("timeout")).getValue();
                            property.timestamp = element.getAttributeByName(new QName("timestamp")).getValue();
                            property.type = switch (element.getName().getLocalPart()) {
                                case "defNumberVector" -> PropertyType.NUMBER;
                                case "defSwitchVector" -> PropertyType.SWITCH;
                                case "defTextVector" -> PropertyType.TEXT;
                                default -> PropertyType.UNKNOWN;
                            };

                            break;
                        case "defNumber":
                        case "defSwitch":
                        case "defText":
                            var name = element.getAttributeByName(new QName("name")).getValue();
                            var label = element.getAttributeByName(new QName("label")).getValue();

                            event = reader.nextEvent();
                            var value = event.asCharacters().getData().trim();

                            if (property != null) {
                                property.values.add(new Value(name, label, value));
                            }
                            break;
                        default:
                            System.out.println("Unhandled element: " + element.getName().getLocalPart());
                    }
                } else if (event.isEndElement()) {
                    var element = event.asEndElement();
                    switch (element.getName().getLocalPart()) {
                        case "defNumberVector":
                        case "defSwitchVector":
                        case "defTextVector":
                            properties.add(property);
                            break;
                    }
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
