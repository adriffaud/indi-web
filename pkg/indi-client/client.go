package indiclient

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Client struct {
	Conn net.Conn
}

func recv(c net.Conn, msg chan string) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		msg <- input.Text()
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
			log.Printf("Received: %s", str)
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
