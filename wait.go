package util

import (
	"context"
	"fmt"
	"time"
)

// waitUntil is a helper function that implements the common waiting pattern
func waitUntil(ctx context.Context, interval time.Duration, maxTries uint, condition func() (bool, error)) error {
	if maxTries == 0 {
		maxTries = 1
	}

	// Try once immediately
	success, err := condition()
	if err != nil {
		return err
	}
	if success {
		return nil
	}

	var tries uint
	for tries = 1; tries < maxTries; tries++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
			success, err = condition()
			if err != nil {
				return err
			}
			if success {
				return nil
			}
		}
	}

	return fmt.Errorf("condition not met after %d tries", maxTries)
}

// WaitFor waits for a function to return true, it will check every interval seconds up until max seconds.
func WaitFor(ctx context.Context, interval time.Duration, maxTries uint, op func() bool) error {
	return waitUntil(ctx, interval, maxTries, func() (bool, error) {
		return op(), nil
	})
}

// WaitForNilError waits for a function to return a nil error, it will check every interval seconds up until max seconds.
func WaitForNilError(ctx context.Context, interval time.Duration, maxTries uint, op func() error) error {
	return waitUntil(ctx, interval, maxTries, func() (bool, error) {
		err := op()
		if err != nil {
			return false, nil // Continue waiting, no error to propagate
		}
		return true, nil
	})
}

// WaitForReturn waits for a function to return a non-nil value, it will check every interval seconds up until max seconds.
// The function returns the value and error returned by the function.
// If maxTries is 0, it will only try once (it will set maxTries internally to 1).
func WaitForReturn[T any](ctx context.Context, interval time.Duration, maxTries uint, op func() (*T, error)) (*T, error) {
	var result *T

	err := waitUntil(ctx, interval, maxTries, func() (bool, error) {
		var err error
		result, err = op()
		if err != nil {
			return false, nil // Continue waiting, don't propagate the error yet
		}
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
