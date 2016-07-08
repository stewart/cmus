package cmus

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn net.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect() error {
	path, err := socketPath()
	if err != nil {
		return err
	}

	addr, err := net.ResolveUnixAddr("unix", path)
	if err != nil {
		return err
	}

	conn, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

// sends a command to cmus
func (c *Client) write(str string) {
	fmt.Fprintf(c.conn, str+"\n")
}

// reads a response from cmus
func (c *Client) read() string {
	lines := []string{}

	scanner := bufio.NewScanner(c.conn)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}

		lines = append(lines, text)
	}

	return strings.Join(lines, "\n")
}

func (c *Client) Status() string {
	c.write("status")
	return c.read()
}
