package cmus

import "testing"

const validStatus = `
status playing
file The Naked and Famous/Passive Me Aggressive You/07 Young Blood.m4a
duration 246
position 13
tag artist The Naked and Famous
tag albumartist The Naked and Famous
tag album Passive Me Aggressive You
tag title Young Blood
tag date 2010
tag compilation no
tag tracknumber 7
set repeat true
set repeat_current false
set shuffle false
set vol_left 100
set vol_right 100
`

func TestParseStatus(t *testing.T) {
	expected := &Status{
		Playing:  true,
		File:     "The Naked and Famous/Passive Me Aggressive You/07 Young Blood.m4a",
		Duration: 246,
		Position: 13,
		Tags: map[string]string{
			"artist":      "The Naked and Famous",
			"albumartist": "The Naked and Famous",
			"album":       "Passive Me Aggressive You",
			"title":       "Young Blood",
			"date":        "2010",
			"compilation": "no",
			"tracknumber": "7",
		},
		Settings: map[string]string{
			"repeat":         "true",
			"repeat_current": "false",
			"shuffle":        "false",
			"vol_left":       "100",
			"vol_right":      "100",
		},
	}

	got, err := parseStatus(validStatus)
	if err != nil {
		t.Error(err)
	}

	if got.Playing != expected.Playing {
		msg := "Expected s.Playing to be %v, got %v"
		t.Errorf(msg, expected.Playing, got.Playing)
	}

	if got.File != expected.File {
		msg := "Expected s.File to be %v, got %v"
		t.Errorf(msg, expected.File, got.File)
	}

	if got.Duration != expected.Duration {
		msg := "Expected s.File to be %v, got %v"
		t.Errorf(msg, expected.Duration, got.Duration)
	}

	if got.Position != expected.Position {
		msg := "Expected s.File to be %v, got %v"
		t.Errorf(msg, expected.Position, got.Position)
	}

	for key, value := range expected.Tags {
		actual := got.Tags[key]

		if actual != value {
			msg := "Expected s.Tags[%q] to be %q, got %q"
			t.Errorf(msg, key, value, actual)
		}
	}

	for key, value := range expected.Settings {
		actual := got.Settings[key]

		if actual != value {
			msg := "Expected s.Settings[%q] to be %q, got %q"
			t.Errorf(msg, key, value, actual)
		}
	}
}

func BenchmarkParseStatus(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			parseStatus(validStatus)
		}
	})
}

func TestStatusTime(t *testing.T) {
	cases := []struct {
		Position, Duration int
		Expected           string
	}{
		{10, 60, "00:10 / 01:00"},
		{63, 90, "01:03 / 01:30"},
		{100, 10000, "01:40 / 166:40"},
	}

	for _, c := range cases {
		status := &Status{Position: c.Position, Duration: c.Duration}
		got := status.Time()

		if got != c.Expected {
			msg := "Expected (&Status{Position: %d, Duration: %d}).Time() to be %q, got %q"
			t.Errorf(msg, c.Position, c.Duration, c.Expected, got)
		}
	}
}
