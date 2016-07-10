package cmus

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

// used to temporarily overwrite ENV vars
type TempEnv struct {
	key, original string
}

func NewTempEnv(key, value string) *TempEnv {
	tmp := &TempEnv{key: key}

	tmp.original = os.Getenv(key)
	os.Setenv(key, value)

	return tmp
}

func (t *TempEnv) Restore() {
	os.Setenv(t.key, t.original)
}

func TestSocketPathWithSocketSet(t *testing.T) {
	expected := "/path/to/socket"

	tmp := NewTempEnv("CMUS_SOCKET", expected)
	defer tmp.Restore()

	actual, err := socketPath()
	if err != nil {
		t.Error(err)
	}

	if expected != actual {
		t.Errorf("Expected socketPath() == %v, got %v", expected, actual)
	}
}

func TestSocketPathWithNoHome(t *testing.T) {
	tmp := NewTempEnv("HOME", "")
	defer tmp.Restore()

	_, err := socketPath()
	if err == nil {
		t.Error("Expected socketPath() to fail, but got no error")
	}

	expected := "environment variable $HOME not set"
	actual := err.Error()

	if actual != expected {
		msg := "Expected socketPath() to error with %q, got %q"
		t.Errorf(msg, expected, actual)
	}
}

func TestExists(t *testing.T) {
	existingFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}

	goneFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}

	os.Remove(goneFile.Name())
	defer os.Remove(existingFile.Name())

	cases := []struct {
		input    string
		expected bool
	}{
		{existingFile.Name(), true},
		{goneFile.Name(), false},
	}

	for _, c := range cases {
		got, err := exists(c.input)
		if err != nil {
			t.Error(err)
		}

		if got != c.expected {
			msg := "Expected exists(%v) == %v, got %v"
			t.Errorf(msg, c.input, c.expected, got)
		}
	}
}

func TestDiscardResult(t *testing.T) {
	testErr := errors.New("Test Error")

	cases := []struct {
		input, expected error
	}{
		{nil, nil},
		{testErr, testErr},
	}

	for _, c := range cases {
		got := discardResult("result", c.input)

		if got != c.expected {
			msg := "Expected discardResult(%v) == %v, got %v"
			t.Errorf(msg, c.input, c.expected, got)
		}
	}
}
