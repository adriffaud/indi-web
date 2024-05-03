package indiclient

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
)

type Property struct {
}

type Client struct {
	conn    net.Conn
	Devices map[string]map[string]any
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:    conn,
		Devices: make(map[string]map[string]any),
	}
	go client.listen(conn)

	return client, nil
}

func (c *Client) GetProperties() error {
	return c.sendMessage("<getProperties version=\"1.7\"/>")
}

func (c *Client) sendMessage(message string) error {
	_, err := fmt.Fprint(c.conn, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (c *Client) listen(conn net.Conn) {
	raw := xml.NewDecoder(conn)
	decoder := xml.NewTokenDecoder(Trimmer{raw})

	for {
		t, err := decoder.Token()
		if t == nil {
			if err == nil {
				continue
			}
			if err == io.EOF {
				log.Println("EOF")
				break
			}
		}

		switch se := t.(type) {
		case xml.StartElement:
			switch se.Name.Local {
			case "defNumberVector":
				var defNumberVector DefNumberVector
				decoder.DecodeElement(&defNumberVector, &se)
				c.addDeviceProperties(defNumberVector)
			case "defSwitchVector":
				var defSwitchVector DefSwitchVector
				decoder.DecodeElement(&defSwitchVector, &se)
				c.addDeviceProperties(defSwitchVector)
			case "defTextVector":
				var defTextVector DefTextVector
				decoder.DecodeElement(&defTextVector, &se)
				c.addDeviceProperties(defTextVector)
			case "defNumber":
				var defNumber DefNumber
				decoder.DecodeElement(&defNumber, &se)
			case "defSwitch":
				var defSwitch DefSwitch
				decoder.DecodeElement(&defSwitch, &se)
			case "defText":
				var defText DefText
				decoder.DecodeElement(&defText, &se)
			default:
				log.Printf("Unhandled data type: %s\n", se.Name.Local)
			}

			fmt.Println("=======================================================")
			for device, groups := range c.Devices {
				fmt.Printf("Device: %s\n", device)
				for group, properties := range groups {
					fmt.Println("---")
					fmt.Printf("Group: %s\n", group)
					fmt.Printf("%+v\n", properties)
				}
			}
		default:
		}
	}
}

func (c *Client) addDeviceProperties(properties VectorAttrs) {
	if _, ok := c.Devices[properties.Device]; !ok {
		c.Devices[properties.Device] = make(map[string]any)
	}

	c.Devices[properties.Device][properties.Group] = properties
}
