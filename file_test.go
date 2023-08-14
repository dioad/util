package util

import (
	"os"
	"testing"
)

func TestExpandPath(t *testing.T) {
	savedVal := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", savedVal)
	}()

	os.Setenv("HOME", "/home/test")
	path, err := ExpandPath("~")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if path == "" {
		t.Errorf("expected path got empty string")
	}
	if path != "/home/test" {
		t.Errorf("expected '/home/test' got '%s'", path)
	}
}
