package cmus

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Status contains information about the current state of cmus, including the
// file that's currently playing, ID3 metadata, and settings.
type Status struct {
	Playing  bool
	File     string
	Duration int
	Position int
	Tags     map[string]string
	Settings map[string]string
}

// Converts the Status' Position/Duration into a displayable playback time
// e.g. "01:30 / 02:00"
func (s *Status) Time() string {
	p := s.Position
	d := s.Duration
	return fmt.Sprintf("%02d:%02d / %02d:%02d", p/60, p%60, d/60, d%60)
}

func parseStatus(input string) (*Status, error) {
	s := &Status{
		Tags:     map[string]string{},
		Settings: map[string]string{},
	}

	lines := strings.Split(strings.TrimSpace(input), "\n")

	for _, line := range lines {
		var a, b, c string
		tokens := strings.SplitN(line, " ", 3)

		if len(tokens) < 2 {
			return s, errors.New("unable to parse cmus status message")
		}

		a = tokens[0]
		b = tokens[1]

		if a == "tag" || a == "set" {
			if len(tokens) < 3 {
				return s, errors.New("unable to parse cmus status message")
			}

			c = tokens[2]
		} else {
			if len(tokens) == 3 {
				b += " " + tokens[2]
			}
		}

		switch a {
		case "status":
			s.Playing = (b == "playing")
		case "file":
			s.File = b
		case "duration":
			n, err := strconv.Atoi(b)
			if err != nil {
				return s, errors.New("unable to parse cmus duration to integer")
			}
			s.Duration = n
		case "position":
			n, err := strconv.Atoi(b)
			if err != nil {
				return s, errors.New("unable to parse cmus position to integer")
			}
			s.Position = n
		case "tag":
			s.Tags[b] = c
		case "set":
			s.Settings[b] = c
		}
	}

	return s, nil
}
