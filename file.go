package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"

	"github.com/dioad/generics"
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

// WaitForFiles waits for a file to exist, it will check every interval seconds up until max seconds.
func WaitForFiles(interval, max int, files ...string) error {
	if interval <= 0 {
		interval = 0
	}
	if max <= 0 {
		max = 1
	}
	for i := 0; i < max; i++ {
		if FilesExist(files...) {
			return nil
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return fmt.Errorf("one or more of %s not found", strings.Join(files, ", "))
}

func fileExists(filename string) error {
	_, err := os.Stat(filename)
	return err
}

// FilesExist checks if all file names exist.
func FilesExist(files ...string) bool {
	return generics.Apply(fileExists, files) == nil
}
