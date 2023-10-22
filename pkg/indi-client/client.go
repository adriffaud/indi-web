package indiclient

import (
	"fmt"
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

	return &Client{Conn: conn}, nil
}

func (c *Client) sendMessage(message string) error {
	_, err := fmt.Fprint(c.Conn, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}

func (c *Client) receiveMessage() (string, error) {
	buffer := make([]byte, 1024)
	_, err := c.Conn.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return string(buffer), nil
}

func (c *Client) GetProperties() (string, error) {
	err := c.sendMessage("<getProperties version=\"1.7\"/>")
	if err != nil {
		return "", fmt.Errorf("failed to send GetProperties request: %v", err)
	}
	response, err := c.receiveMessage()
	if err != nil {
		return "", fmt.Errorf("failed to receive GetProperties response: %v", err)
	}

	return response, nil
}
