package indiclient

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
)

type Client struct {
	conn    net.Conn
	Channel chan any
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	ch := make(chan any)
	go recv(conn, ch)

	return &Client{
		conn:    conn,
		Channel: ch,
	}, nil
}

func (c *Client) sendMessage(message string) error {
	_, err := fmt.Fprint(c.conn, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (c *Client) GetProperties() error {
	return c.sendMessage("<getProperties version=\"1.7\"/>")
}

// Trimmer is used to remove blank space from received XML
type Trimmer struct {
	dec *xml.Decoder
}

func (tr Trimmer) Token() (xml.Token, error) {
	t, err := tr.dec.Token()
	if cd, ok := t.(xml.CharData); ok {
		t = xml.CharData(bytes.TrimSpace(cd))
	}
	return t, err
}

func recv(c net.Conn, ch chan<- any) {
	raw := xml.NewDecoder(c)
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
				ch <- defNumberVector
			case "defSwitchVector":
				var defSwitchVector DefSwitchVector
				decoder.DecodeElement(&defSwitchVector, &se)
				ch <- defSwitchVector
			case "defTextVector":
				var defTextVector DefTextVector
				decoder.DecodeElement(&defTextVector, &se)
				ch <- defTextVector
			case "defNumber":
				var defNumber DefNumber
				decoder.DecodeElement(&defNumber, &se)
				ch <- defNumber
			case "defSwitch":
				var defSwitch DefSwitch
				decoder.DecodeElement(&defSwitch, &se)
				ch <- defSwitch
			case "defText":
				var defText DefText
				decoder.DecodeElement(&defText, &se)
				ch <- defText
			default:
				log.Printf("Unhandled data type: %s\n", se.Name.Local)
			}
		default:
		}
	}
}
