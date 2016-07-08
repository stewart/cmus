package cmus

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"sync"
)

type Client struct {
	conn net.Conn
	mut  sync.Mutex
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
func (c *Client) write(str string) error {
	if c.conn == nil {
		return errors.New("client is not connected")
	}

	_, err := c.conn.Write([]byte(str + "\n"))
	if err != nil {
		return err
	}

	return nil
}

// reads a response from cmus
func (c *Client) read() (string, error) {
	lines := []string{}

	scanner := bufio.NewScanner(c.conn)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}

		lines = append(lines, text)
	}

	return strings.Join(lines, "\n"), scanner.Err()
}

// executes a command against cmus, and returns the result
func (c *Client) Cmd(command string) (string, error) {
	c.mut.Lock()
	defer c.mut.Unlock()

	if err := c.write(command); err != nil {
		return "", err
	}

	return c.read()
}

func (c *Client) Status() (string, error) {
	return c.Cmd("status")
}

func (c *Client) Play(filename string) (string, error) {
	return c.Cmd("player-play " + filename)
}

// Toggles play/pause
func (c *Client) PlayPause() (string, error) {
	return c.Cmd("player-pause")
}

func (c *Client) Stop() (string, error) {
	return c.Cmd("player-stop")
}

func (c *Client) Prev() (string, error) {
	return c.Cmd("player-prev")
}

func (c *Client) Next() (string, error) {
	return c.Cmd("player-prev")
}

func (c *Client) Seek(time string) (string, error) {
	return c.Cmd("seek " + time)
}

func (c *Client) Volume(level string) (string, error) {
	return c.Cmd("vol " + level)
}

func (c *Client) Shuffle() (string, error) {
	return c.Cmd("toggle shuffle")
}

func (c *Client) Repeat() (string, error) {
	return c.Cmd("toggle repeat")
}
