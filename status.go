package cmus

import (
	"errors"
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

func parseStatus(input string) (*Status, error) {
	s := &Status{
		Tags:     map[string]string{},
		Settings: map[string]string{},
	}

	lines := strings.Split(strings.TrimSpace(input), "\n")

	for _, line := range lines {
		var a, b, c string
		tokens := strings.SplitN(line, " ", 2)

		if len(tokens) < 2 {
			return s, errors.New("unable to parse cmus status message")
		}

		a = tokens[0]

		if a == "tag" || a == "set" {
			tokens = strings.SplitN(line, " ", 3)

			if len(tokens) < 3 {
				return s, errors.New("unable to parse cmus status message")
			}

			c = tokens[2]
		}

		b = tokens[1]

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
