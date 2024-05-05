package indiclient

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
)

type PropertyType int64

const (
	Number PropertyType = iota
	Switch
	Text
)

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
	Format    string
	Min       string
	Max       string
	Step      string
	Rule      string
	Values    []interface{}
}

type Client struct {
	conn    net.Conn
	Devices map[string]map[string][]Property
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:    conn,
		Devices: make(map[string]map[string][]Property),
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

				property := Property{
					Device: defNumberVector.Device,
					Group:  defNumberVector.Group,
					Type:   Number,
					Name:   defNumberVector.Label,
				}

				c.addToTree(property)
			case "defSwitchVector":
				var defSwitchVector DefSwitchVector
				decoder.DecodeElement(&defSwitchVector, &se)

				property := Property{
					Device: defSwitchVector.Device,
					Group:  defSwitchVector.Group,
					Type:   Switch,
					Name:   defSwitchVector.Label,
				}

				c.addToTree(property)
			case "defTextVector":
				var defTextVector DefTextVector
				decoder.DecodeElement(&defTextVector, &se)

				property := Property{
					Device: defTextVector.Device,
					Group:  defTextVector.Group,
					Type:   Text,
					Name:   defTextVector.Label,
				}

				c.addToTree(property)
			default:
				log.Printf("!!!! Unhandled data type: %s\n", se.Name.Local)
			}

			fmt.Println("=======================================================")
			for device, groups := range c.Devices {
				fmt.Printf("Device: %s\n", device)
				for group, properties := range groups {
					fmt.Println("---")
					fmt.Printf("Group: %s\n", group)
					for _, property := range properties {
						fmt.Printf("%+v\n", property)
					}
				}
			}
		default:
		}
	}
}

func (c *Client) addToTree(property Property) {
	if _, exists := c.Devices[property.Device]; !exists {
		c.Devices[property.Device] = make(map[string][]Property)
	}
	if _, exists := c.Devices[property.Device][property.Group]; !exists {
		c.Devices[property.Device][property.Group] = make([]Property, 0)
	}
	c.Devices[property.Device][property.Group] = append(c.Devices[property.Device][property.Group], property)
}
