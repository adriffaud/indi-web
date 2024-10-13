package indiclient

import (
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net"
	"slices"
	"strconv"
	"sync"
	"time"
)

const timeout = 5 * time.Second

type Client struct {
	conn       net.Conn
	Properties Properties
	observers  map[Observer]struct{}
	mu         sync.Mutex
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:       conn,
		Properties: make([]Property, 0),
		observers:  map[Observer]struct{}{},
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

func (c *Client) Connect(driver string) error {
	return c.NewPropertyValue(PropertySelector{Device: driver, Name: "CONNECTION", ValueName: "CONNECT"})
}

func (c *Client) NewPropertyValue(selector PropertySelector) error {
	property := c.Properties.FindProperty(selector)

	var newValue string
	for _, value := range property.Values {
		if value.Name == selector.ValueName && value.Value == "Off" {
			newValue = "On"
		} else if value.Name == selector.ValueName && value.Value == "On" {
			newValue = "Off"
		}
	}

	xml := fmt.Sprintf("<newSwitchVector device=\"%s\" name=\"%s\"><oneSwitch name=\"%s\">%s</oneSwitch></newSwitchVector>", selector.Device, selector.Name, selector.ValueName, newValue)

	slog.Debug("sending new property value", "selector", selector, "xml", xml)

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
	inactivityTimer := time.NewTimer(timeout)

	go func() {
		for {
			<-inactivityTimer.C
			slog.Debug("😴😴 connection idle")
			c.Notify(Event{EventType: Timeout})
		}
	}()

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

		inactivityTimer.Reset(timeout)

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
				for _, attr := range se.Attr {
					if attr.Name.Local == "message" {
						c.Notify(Event{EventType: Message, Message: attr.Value})
					}
				}
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
			slog.Warn(fmt.Sprintf("⚠️⚠️⚠️ Unhandled element type: %T\n", t), "value", se)
		}
	}

	inactivityTimer.Stop()
	slog.Debug("connection closed")
}

func (c *Client) addToProperties(property Property) {
	c.delFromProperties(property.Device, property.Name)
	slog.Debug("adding property", "property", property)
	c.Properties = append(c.Properties, property)
	c.Notify(Event{EventType: Add, Property: property})
}

func (c *Client) delFromProperties(device, name string) {
	slog.Debug("deleting property", "device", device, "name", name)
	propIdx := slices.IndexFunc(c.Properties, func(p Property) bool { return p.Device == device && p.Name == name })

	if propIdx < 0 {
		return
	}

	prop := c.Properties.FindProperty(PropertySelector{Device: device, Name: name})
	c.Notify(Event{EventType: Delete, Property: *prop})

	c.Properties = append(c.Properties[:propIdx], c.Properties[propIdx+1:]...)
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

	c.Notify(Event{EventType: Update, Property: *prop})
}

func (c *Client) Register(o Observer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.observers[o] = struct{}{}
	slog.Debug("➕ Adding observer", "count", len(c.observers))
}

func (c *Client) Unregister(o Observer) {
	delete(c.observers, o)
	slog.Debug("➖ Removing observer", "count", len(c.observers))
}

func (c *Client) Notify(e Event) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for o := range c.observers {
		o.OnNotify(e)
	}
}
