package indiclient

import (
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net"
	"slices"
	"strconv"
)

type Client struct {
	conn       net.Conn
	Properties Properties
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:       conn,
		Properties: make([]Property, 0),
	}

	go client.listen(conn)

	return client, nil
}

func (c *Client) Close() {
	c.conn.Close()
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
				slog.Info("EOF")
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
					Device:    defNumberVector.Device,
					Group:     defNumberVector.Group,
					Type:      Number,
					Name:      defNumberVector.Name,
					Label:     defNumberVector.Label,
					State:     defNumberVector.State,
					Perm:      defNumberVector.Perm,
					Timeout:   defNumberVector.Timeout,
					Timestamp: defNumberVector.Timestamp,
				}

				values := make([]Value, 0)
				for _, number := range defNumberVector.DefNumber {
					values = append(values, Value{Name: number.Name, Label: number.Label, Value: strconv.Itoa(number.Value)})
				}
				property.Values = values

				c.addToProperties(property)
			case "defSwitchVector":
				var defSwitchVector DefSwitchVector
				decoder.DecodeElement(&defSwitchVector, &se)

				property := Property{
					Device:    defSwitchVector.Device,
					Group:     defSwitchVector.Group,
					Type:      Switch,
					Name:      defSwitchVector.Name,
					Label:     defSwitchVector.Label,
					State:     defSwitchVector.State,
					Perm:      defSwitchVector.Perm,
					Timeout:   defSwitchVector.Timeout,
					Timestamp: defSwitchVector.Timestamp,
					Rule:      defSwitchVector.Rule,
				}

				values := make([]Value, 0)
				for _, item := range defSwitchVector.DefSwitch {
					values = append(values, Value{Name: item.Name, Label: item.Label, Value: item.Value})
				}
				property.Values = values

				c.addToProperties(property)
			case "defTextVector":
				var defTextVector DefTextVector
				decoder.DecodeElement(&defTextVector, &se)

				property := Property{
					Device:    defTextVector.Device,
					Group:     defTextVector.Group,
					Type:      Text,
					Name:      defTextVector.Name,
					Label:     defTextVector.Label,
					State:     defTextVector.State,
					Perm:      defTextVector.Perm,
					Timeout:   defTextVector.Timeout,
					Timestamp: defTextVector.Timestamp,
				}

				values := make([]Value, 0)
				for _, text := range defTextVector.DefText {
					values = append(values, Value{Name: text.Name, Label: text.Label, Value: text.Value})
				}
				property.Values = values

				c.addToProperties(property)
			default:
				// slog.Warn("Unhandled data type", "type", se.Name.Local)
			}
		default:
			// slog.Warn(fmt.Sprintf("Unhandled element type: %T\n", t))
		}
	}
}

func (c *Client) addToProperties(property Property) {
	var containsFunc = func(p Property) bool {
		return p.Device == property.Device && p.Group == property.Group && p.Name == property.Name
	}

	if !slices.ContainsFunc(c.Properties, containsFunc) {
		c.Properties = append(c.Properties, property)
	}
}
