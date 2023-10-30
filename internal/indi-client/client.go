package indiclient

import (
	"encoding/xml"
	"fmt"
	"log"
	"net"
)

type Client struct {
	Conn net.Conn
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	go recv(conn)

	return &Client{Conn: conn}, nil
}

func (c *Client) sendMessage(message string) error {
	_, err := fmt.Fprint(c.Conn, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (c *Client) GetProperties() error {
	return c.sendMessage("<getProperties version=\"1.7\"/>")
}

func recv(c net.Conn) {
	decoder := xml.NewDecoder(c)

	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}

		switch t.(type) {
		case xml.StartElement:
			log.Println("============================")
			log.Printf("%+v\n", t)
		case xml.EndElement:
		}
	}
}
