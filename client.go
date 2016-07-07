package cmus

import (
	"fmt"
	"io"
	"net"
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
func (c *Client) read() (string, error) {
	var data string

	buf := make([]byte, 8192)
	newline := false
	len := 0

	for {
		rc, err := c.conn.Read(buf)
		if err != nil && err != io.EOF {
			return data, err
		}

		len += rc

		if newline && buf[0] == '\n' {
			break
		}

		if len == 1 && buf[0] == '\n' {
			break
		}

		if rc > 1 && buf[rc-1] == '\n' && buf[rc-2] == '\n' {
			data += string(buf[:rc-1])
			break
		}

		newline = buf[rc-1] == '\n'
		data += string(buf)
	}

	return data, nil
}

func (c *Client) Status() {
}
