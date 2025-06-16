package util

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWaitFor(t *testing.T) {
	t.Run("immediate success", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0

		err := WaitFor(ctx, 10*time.Millisecond, 5, func() bool {
			callCount++
			return true
		})

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if callCount != 1 {
			t.Errorf("expected 1 call, got: %d", callCount)
		}
	})

	t.Run("eventual success", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0

		err := WaitFor(ctx, 10*time.Millisecond, 5, func() bool {
			callCount++
			return callCount == 3
		})

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if callCount != 3 {
			t.Errorf("expected 3 calls, got: %d", callCount)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0

		err := WaitFor(ctx, 10*time.Millisecond, 3, func() bool {
			callCount++
			return false
		})

		if err == nil {
			t.Error("expected error, got nil")
		}
		if callCount != 3 {
			t.Errorf("expected 3 calls, got: %d", callCount)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		callCount := 0

		// Cancel after the first call
		go func() {
			time.Sleep(15 * time.Millisecond)
			cancel()
		}()

		err := WaitFor(ctx, 10*time.Millisecond, 5, func() bool {
			callCount++
			return false
		})

		if err == nil {
			t.Error("expected error, got nil")
		}
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled error, got: %v", err)
		}
	})

	t.Run("zero maxTries", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0

		err := WaitFor(ctx, 10*time.Millisecond, 0, func() bool {
			callCount++
			return false
		})

		if err == nil {
			t.Error("expected error, got nil")
		}
		if callCount != 1 {
			t.Errorf("expected 1 call, got: %d", callCount)
		}
	})
}

func TestWaitForNilError(t *testing.T) {
	t.Run("immediate success", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0

		err := WaitForNilError(ctx, 10*time.Millisecond, 5, func() error {
			callCount++
			return nil
		})

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if callCount != 1 {
			t.Errorf("expected 1 call, got: %d", callCount)
		}
	})

	t.Run("eventual success", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0

		err := WaitForNilError(ctx, 10*time.Millisecond, 5, func() error {
			callCount++
			if callCount < 3 {
				return errors.New("not ready yet")
			}
			return nil
		})

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if callCount != 3 {
			t.Errorf("expected 3 calls, got: %d", callCount)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0

		err := WaitForNilError(ctx, 10*time.Millisecond, 3, func() error {
			callCount++
			return errors.New("always failing")
		})

		if err == nil {
			t.Error("expected error, got nil")
		}
		if callCount != 3 {
			t.Errorf("expected 3 calls, got: %d", callCount)
		}
	})
}

func TestWaitForReturn(t *testing.T) {
	t.Run("immediate success", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0
		expectedResult := "success"

		result, err := WaitForReturn(ctx, 10*time.Millisecond, 5, func() (*string, error) {
			callCount++
			return &expectedResult, nil
		})

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if callCount != 1 {
			t.Errorf("expected 1 call, got: %d", callCount)
		}
		if result == nil || *result != expectedResult {
			t.Errorf("expected result %v, got: %v", expectedResult, result)
		}
	})

	t.Run("eventual success", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0
		expectedResult := "success"

		result, err := WaitForReturn(ctx, 10*time.Millisecond, 5, func() (*string, error) {
			callCount++
			if callCount < 3 {
				return nil, errors.New("not ready yet")
			}
			return &expectedResult, nil
		})

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
		if callCount != 3 {
			t.Errorf("expected 3 calls, got: %d", callCount)
		}
		if result == nil || *result != expectedResult {
			t.Errorf("expected result %v, got: %v", expectedResult, result)
		}
	})

	t.Run("nil result", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0

		result, err := WaitForReturn(ctx, 10*time.Millisecond, 3, func() (*string, error) {
			callCount++
			return nil, nil
		})

		if err == nil {
			t.Error("expected error, got nil")
		}
		if callCount != 3 {
			t.Errorf("expected 3 calls, got: %d", callCount)
		}
		if result != nil {
			t.Errorf("expected nil result, got: %v", result)
		}
	})

	t.Run("error result", func(t *testing.T) {
		ctx := context.Background()
		callCount := 0
		expectedError := errors.New("test error")

		result, err := WaitForReturn(ctx, 10*time.Millisecond, 3, func() (*string, error) {
			callCount++
			if callCount == 3 {
				return nil, expectedError
			}
			return nil, nil
		})

		if err == nil {
			t.Error("expected error, got nil")
		}
		if callCount != 3 {
			t.Errorf("expected 3 calls, got: %d", callCount)
		}
		if result != nil {
			t.Errorf("expected nil result, got: %v", result)
		}
	})
}

// TestWaitForFiles tests the WaitForFiles function which uses WaitFor internally
func TestWaitForFiles(t *testing.T) {
	t.Run("files exist", func(t *testing.T) {
		// Create a temporary file
		tmpfile, err := os.CreateTemp("", "example")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		ctx := context.Background()
		err = WaitForFiles(ctx, 1, 2, tmpfile.Name())

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("files don't exist", func(t *testing.T) {
		ctx := context.Background()
		err := WaitForFiles(ctx, 1, 2, "/path/to/nonexistent/file")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// TestWaitForFile tests the WaitForFile function which uses WaitFor internally
func TestWaitForFile(t *testing.T) {
	t.Run("file exists", func(t *testing.T) {
		// Create a temporary file
		tmpfile, err := os.CreateTemp("", "example")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpfile.Name())

		ctx := context.Background()
		err = WaitForFile(ctx, 10*time.Millisecond, 2, tmpfile.Name())

		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("file doesn't exist", func(t *testing.T) {
		ctx := context.Background()
		err := WaitForFile(ctx, 10*time.Millisecond, 2, "/path/to/nonexistent/file")

		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("file created during wait", func(t *testing.T) {
		// Create a temporary directory
		tmpdir, err := os.MkdirTemp("", "example")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tmpdir)

		// Path to a file that doesn't exist yet
		filePath := filepath.Join(tmpdir, "newfile.txt")

		// Start waiting for the file in a goroutine
		ctx := context.Background()
		errChan := make(chan error)
		go func() {
			errChan <- WaitForFile(ctx, 10*time.Millisecond, 10, filePath)
		}()

		// Wait a bit and then create the file
		time.Sleep(30 * time.Millisecond)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()

		// Check that the wait function returns without error
		select {
		case err := <-errChan:
			if err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		case <-time.After(200 * time.Millisecond):
			t.Error("timeout waiting for WaitForFile to return")
		}
	})
}
