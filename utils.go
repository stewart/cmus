package cmus

import (
	"errors"
	"os"
	"path/filepath"
)

// Finds the path to cmus' unix socket
func socketPath() (string, error) {
	socket := os.Getenv("CMUS_SOCKET")
	if socket != "" {
		return socket, nil
	}

	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("environment variable $HOME not set")
	}

	xdg := os.Getenv("XDG_RUNTIME_DIR")

	config_path := filepath.Join(home, ".cmus")
	config_path_exists, err := exists(config_path)
	if err != nil {
		return "", err
	}

	if config_path_exists {
		return config_path, nil
	} else if xdg == "" {
		return filepath.Join(home, ".config", "cmus", "socket"), nil
	} else {
		return filepath.Join(xdg, "cmus-socket"), nil
	}
}

// checks if a file path exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}
