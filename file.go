package util

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

func CleanOpen(path string) (*os.File, error) {
	path, err := ExpandPath(path)
	if err != nil {
		return nil, err
	}

	path = filepath.Clean(path)

	return os.Open(path)
}

func CleanOpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	expandedPath, err := ExpandPath(path)
	if err != nil {
		return nil, err
	}

	cleanPath := filepath.Clean(expandedPath)

	return os.OpenFile(cleanPath, flag, perm) // #nosec
}

// CreateDirPath creates a directory path if it doesn't exist.
func CreateDirPath(path string, defaultPath string) (string, error) {
	if path == "" {
		path = defaultPath
	}

	path, err := ExpandPath(path)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(path, 0750)
	if err != nil {
		return "", err
	}

	return path, nil
}

// ExpandPath expands a path to an absolute path.
// It also expands ~ and environment variables.
func ExpandPath(path string) (string, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return "", err
	}

	path = os.ExpandEnv(path)

	path = filepath.Clean(path)

	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return path, nil
}
