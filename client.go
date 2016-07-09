package cmus

import (
	"bufio"
	"errors"
	"net"
	"regexp"
	"strings"
	"sync"
)

var errorRegexp = regexp.MustCompile(`^Error\:\s+(.+)\s*$`)

// Client is a connection to a cmus server.
type Client struct {
	conn net.Conn
	mut  sync.Mutex
}

// Connects the Client to cmus. The path socket is automatically determined,
// based on the same logic cmus uses internally.
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

// Sends a command (provided as a string) to cmus. Will return an error if the
// client is not connected to cmus, or if an error occured while writing.
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

// Reads a response from cmus. Will block until one is received. Returns an
// error if one occured during response fetching.
func (c *Client) read() (string, error) {
	lines := []string{}

	scanner := bufio.NewScanner(c.conn)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}

		// check for error message
		errs := errorRegexp.FindStringSubmatch(text)
		if errs != nil {
			return "", errors.New("cmus error: " + errs[1])
		}

		lines = append(lines, text)
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

// Cmd runs a command against cmus, and returns the result of that command.
func (c *Client) Cmd(command string) (string, error) {
	c.mut.Lock()
	defer c.mut.Unlock()

	if err := c.write(command); err != nil {
		return "", err
	}

	return c.read()
}

// Status is a shorthand for Cmd("status"). It returns status information from
// cmus, including metadata about the currently playing track and various
// settings.
func (c *Client) Status() (string, error) {
	return c.Cmd("status")
}

// Play is a shorthand for Cmd("player-play " + filename). It plays the given
// track, or, if none is specified, [re]plays the current track from the
// beginning.
func (c *Client) Play(filename string) error {
	return discardResult(c.Cmd("player-play " + filename))
}

// PlayPause is a shorthand for Cmd("player-pause"). It toggles pause.
func (c *Client) PlayPause() error {
	return discardResult(c.Cmd("player-pause"))
}

// Stop is a shorthand for Cmd("player-stop"). It stops playback.
func (c *Client) Stop() error {
	return discardResult(c.Cmd("player-stop"))
}

// Prev is a shorthand for Cmd("player-prev"). It skips to the previous track.
func (c *Client) Prev() error {
	return discardResult(c.Cmd("player-prev"))
}

// Next is a shorthand for Cmd("player-prev"). It skips to the next track.
func (c *Client) Next() error {
	return discardResult(c.Cmd("player-prev"))
}

// Seek is a shorthand for Cmd("seek " + time). It seeks to an absolute or
// relative position. Position is given in the format "[+-](<num>mh |
// [HH:]MM:SS)". Some examples:
//
// Seek 1 minute backward
//     client.Seek("-1m")
//
// Seek 5 seconds forward
//     client.Seek("+5")
//
// Seek to absolute position 1h
//     client.Seek("1h")
//
// Seek 90 seconds forward
//     client.Seek("+1:30")
func (c *Client) Seek(time string) error {
	return discardResult(c.Cmd("seek " + time))
}

// Volume is a shorthand for Cmd("vol " + level). It sets, increases, or
// decreases volume.
//
// To increase or decrease volume prefix the value with - or +, otherwise value
// is treated as absolute volume.
//
// Both absolute and relative values can be given as percentage units (suffixed
// with %) or as internal values (hardware may have volume in range 0-31 for
// example).
func (c *Client) Volume(level string) error {
	return discardResult(c.Cmd("vol " + level))
}

// Shuffle is a shorthand for Cmd("toggle shuffle"). It toggles whether or not
// the playback order is shuffled.
func (c *Client) Shuffle() error {
	return discardResult(c.Cmd("toggle shuffle"))
}

// Repeat is a shorthand for Cmd("toggle repeat"). It toggles whether or not
// playback will repeat after all tracks are played.
func (c *Client) Repeat() error {
	return discardResult(c.Cmd("toggle repeat"))
}
