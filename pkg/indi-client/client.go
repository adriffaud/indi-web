package indiclient

import (
	"fmt"
	"log"
	"net"
)

type Client struct {
	Conn net.Conn
}

func recv(c net.Conn, msgch chan string) {
	buf := make([]byte, 2048)
	for {
		n, err := c.Read(buf)
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		msg := buf[:n]
		log.Printf("[INDI Client] Received: %s", msg)
		// msgch <- string(msg)
	}
}

func New(address string) (*Client, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	msg := make(chan string)
	go recv(conn, msg)

	go func() {
		for {
			str := <-msg
			log.Printf("[INDI Client] Received: %s", str)
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
