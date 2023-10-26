package indiclient

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"net"
)

type Client struct {
	Conn net.Conn
}

type Message struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
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
	buf := make([]byte, 2048)
	var incompleteData []byte

	for {
		n, err := c.Read(buf)
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		data := append(incompleteData, buf[:n]...)
		incompleteData = processData(data)
	}
}

func processData(data []byte) []byte {
	decoder := xml.NewDecoder(bytes.NewReader(data))

	for {
		var msg Message
		err := decoder.Decode(&msg)
		if err != nil {
			return data
		}

		log.Println("================================")
		log.Printf("[INDI Client] Received:\n%+v\n", msg)
	}
}
