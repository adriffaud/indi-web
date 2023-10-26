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

	msg := make(chan Message)
	go recv(conn, msg)

	go func() {
		for {
			str := <-msg
			log.Printf("[INDI Client] Received: %s\n", str)
		}
	}()

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

func recv(c net.Conn, msgch chan Message) {
	buf := make([]byte, 2048)
	var incompleteData []byte

	for {
		n, err := c.Read(buf)
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		data := append(incompleteData, buf[:n]...)
		incompleteData = processData(data, msgch)
	}
}

func processData(data []byte, msgch chan Message) []byte {
	decoder := xml.NewDecoder(bytes.NewReader(data))

	for {
		var msg Message
		err := decoder.Decode(&msg)
		if err != nil {
			return data
		}

		msgch <- msg
	}
}
