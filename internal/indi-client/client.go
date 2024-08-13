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

func (c *Client) NewPropertyValue(selector PropertySelector) error {
	property := c.Properties.FindProperty(selector)

	var newValue string
	for _, value := range property.Values {
		if value.Name == selector.Value && value.Value == "Off" {
			newValue = "On"
		} else if value.Name == selector.Value && value.Value == "On" {
			newValue = "Off"
		}
	}

	xml := fmt.Sprintf("<newSwitchVector device=\"%s\" name=\"%s\"><oneSwitch name=\"%s\">%s</oneSwitch></newSwitchVector>", selector.Device, selector.Name, selector.Value, newValue)
	slog.Debug("NewPropertyValue", "selector", selector, "property", property, "xml", xml)

	return c.sendMessage(xml)
}

func (c *Client) sendMessage(message string) error {
	_, err := fmt.Fprint(c.conn, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (c *Client) listen(reader io.Reader) {
	raw := xml.NewDecoder(reader)
	decoder := xml.NewTokenDecoder(Trimmer{raw})

	var property Property
	var value Value

	for {
		t, err := decoder.Token()
		if t == nil {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
		}

		switch se := t.(type) {
		case xml.StartElement:
			attrs := make(map[string]string)
			for _, attr := range se.Attr {
				attrs[attr.Name.Local] = attr.Value
			}

			switch se.Name.Local {
			case "defNumberVector", "defSwitchVector", "defTextVector":
				property = Property{
					Device:    attrs["device"],
					Group:     attrs["group"],
					Name:      attrs["name"],
					Label:     attrs["label"],
					State:     attrs["state"],
					Perm:      attrs["perm"],
					Timestamp: attrs["timestamp"],
					Rule:      attrs["rule"],
				}

				switch se.Name.Local {
				case "defNumberVector":
					property.Type = Number
				case "defSwitchVector":
					property.Type = Switch
				case "defTextVector":
					property.Type = Text
				}

				if timeout, err := strconv.Atoi(attrs["timeout"]); err == nil {
					property.Timeout = timeout
				}
			case "defNumber", "defSwitch", "defText":
				value = Value{
					Name:  attrs["name"],
					Label: attrs["label"],
				}
			case "delProperty":
				c.delFromProperties(attrs["device"], attrs["name"])
			case "setNumberVector", "setSwitchVector":
				property = Property{
					Device:    attrs["device"],
					Name:      attrs["name"],
					State:     attrs["state"],
					Timestamp: attrs["timestamp"],
				}
			case "oneNumber", "oneSwitch":
				attrs := make(map[string]string)
				for _, attr := range se.Attr {
					attrs[attr.Name.Local] = attr.Value
				}
				value = Value{Name: attrs["name"], Label: attrs["label"]}
			case "message":
				slog.Debug("MESSAGE", "message", se)
			default:
				slog.Warn("Unhandled data type", "type", se.Name.Local, "raw", se)
			}
		case xml.CharData:
			value.Value = string(se)
		case xml.EndElement:
			switch se.Name.Local {
			case "defNumberVector", "defSwitchVector", "defTextVector":
				c.addToProperties(property)
			case "defNumber", "defSwitch", "defText", "oneNumber", "oneSwitch":
				property.Values = append(property.Values, value)
			case "setNumberVector", "setSwitchVector":
				c.updatePropertyValues(property)
			}
		default:
			slog.Warn(fmt.Sprintf("⚠️ Unhandled element type: %T\n", t), "value", se)
		}
	}
}

func (c *Client) addToProperties(property Property) {
	c.delFromProperties(property.Device, property.Name)
	slog.Debug("➕ Adding property", "property", property)
	c.Properties = append(c.Properties, property)
}

func (c *Client) delFromProperties(device, name string) {
	slog.Debug("🚮 Deleting property", "device", device, "name", name)
	propIdx := slices.IndexFunc(c.Properties, func(p Property) bool { return p.Device == device && p.Name == name })
	if propIdx >= 0 {
		c.Properties = append(c.Properties[:propIdx], c.Properties[propIdx+1:]...)
	}
}

func (c *Client) updatePropertyValues(property Property) {
	propIdx := slices.IndexFunc(c.Properties, func(p Property) bool { return p.Device == property.Device && p.Name == property.Name })
	if propIdx == -1 {
		panic("trying to update unexisting property")
	}

	prop := &c.Properties[propIdx]
	prop.State = property.State
	prop.Timestamp = property.Timestamp

	for _, newValue := range property.Values {
		oldValueIdx := slices.IndexFunc(prop.Values, func(v Value) bool { return v.Name == newValue.Name })
		prop.Values[oldValueIdx].Value = newValue.Value
	}
}
