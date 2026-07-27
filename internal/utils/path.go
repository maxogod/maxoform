package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func ExpandHome(path string) (string, error) {
	if path == "~" {
		return os.UserHomeDir()
	}
	if !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, strings.TrimPrefix(path, "~/")), nil
}
