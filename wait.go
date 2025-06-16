package util

import (
	"context"
	"fmt"
	"os"
	"time"
)

// waitUntil is a helper function that implements the common waiting pattern.
// It tries the condition immediately and then at regular intervals until it succeeds,
// the context is canceled, or the maximum number of tries is reached.
//
// Parameters:
//   - ctx: Context for cancellation
//   - interval: Time to wait between tries
//   - maxTries: Maximum number of times to try the condition (including the immediate try)
//   - condition: Function that returns (success, error)
//
// Returns:
//   - error: nil if condition succeeded, otherwise an error explaining why it failed
func waitUntil(ctx context.Context, interval time.Duration, maxTries uint, condition func() (bool, error)) error {
	if maxTries == 0 {
		maxTries = 1
	}

	// Try once immediately
	success, err := condition()
	if err != nil {
		return fmt.Errorf("condition failed with error: %w", err)
	}
	if success {
		return nil
	}

	// Use a timer instead of time.After to avoid potential resource leaks
	timer := time.NewTimer(interval)
	defer timer.Stop()

	var tries uint
	for tries = 1; tries < maxTries; tries++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("waiting canceled: %w", ctx.Err())
		case <-timer.C:
			success, err = condition()
			if err != nil {
				return fmt.Errorf("condition failed with error on try %d: %w", tries+1, err)
			}
			if success {
				return nil
			}
			// Reset the timer for the next interval
			timer.Reset(interval)
		}
	}

	return fmt.Errorf("condition not met after %d tries", maxTries)
}

// WaitFor waits for a function to return true.
//
// It will check the function immediately and then every interval duration until:
// - The function returns true
// - The context is canceled
// - The maximum number of tries is reached
//
// Parameters:
//   - ctx: Context for cancellation
//   - interval: Time to wait between tries
//   - maxTries: Maximum number of times to try the condition (including the immediate try)
//   - op: Function that returns true when the condition is met
//
// Returns:
//   - error: nil if condition succeeded, otherwise an error explaining why it failed
//
// Example:
//
//	// Wait for a service to be ready, checking every 2 seconds, up to 30 tries
//	err := util.WaitFor(ctx, 2*time.Second, 30, func() bool {
//	    return isServiceReady()
//	})
func WaitFor(ctx context.Context, interval time.Duration, maxTries uint, op func() bool) error {
	return waitUntil(ctx, interval, maxTries, func() (bool, error) {
		return op(), nil
	})
}

// WaitForNilError waits for a function to return a nil error.
//
// It will check the function immediately and then every interval duration until:
// - The function returns nil error
// - The context is canceled
// - The maximum number of tries is reached
//
// Parameters:
//   - ctx: Context for cancellation
//   - interval: Time to wait between tries
//   - maxTries: Maximum number of times to try the condition (including the immediate try)
//   - op: Function that returns nil error when the condition is met
//
// Returns:
//   - error: nil if condition succeeded, otherwise an error explaining why it failed
//
// Example:
//
//	// Wait for a database connection to be established, checking every second, up to 10 tries
//	err := util.WaitForNilError(ctx, time.Second, 10, func() error {
//	    return db.Ping()
//	})
func WaitForNilError(ctx context.Context, interval time.Duration, maxTries uint, op func() error) error {
	return waitUntil(ctx, interval, maxTries, func() (bool, error) {
		err := op()
		if err != nil {
			return false, nil // Continue waiting, no error to propagate
		}
		return true, nil
	})
}

// WaitForReturn waits for a function to return a non-nil value.
//
// It will check the function immediately and then every interval duration until:
// - The function returns a non-nil value and nil error
// - The context is canceled
// - The maximum number of tries is reached
//
// Parameters:
//   - ctx: Context for cancellation
//   - interval: Time to wait between tries
//   - maxTries: Maximum number of times to try the condition (including the immediate try)
//   - op: Function that returns a non-nil value and nil error when the condition is met
//
// Returns:
//   - *T: The value returned by the function when it succeeds
//   - error: nil if condition succeeded, otherwise an error explaining why it failed
//
// If maxTries is 0, it will be set to 1 (try only once).
//
// Example:
//
//	// Wait for a resource to be created, checking every 5 seconds, up to 20 tries
//	resource, err := util.WaitForReturn(ctx, 5*time.Second, 20, func() (*Resource, error) {
//	    return client.GetResource(resourceID)
//	})
func WaitForReturn[T any](ctx context.Context, interval time.Duration, maxTries uint, op func() (*T, error)) (*T, error) {
	var result *T

	err := waitUntil(ctx, interval, maxTries, func() (bool, error) {
		var err error
		result, err = op()
		if err != nil {
			return false, nil // Continue waiting, don't propagate the error yet
		}
		if result == nil {
			return false, nil // Continue waiting, we need a non-nil result
		}
		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get result: %w", err)
	}

	return result, nil
}

// WaitForFile waits for a file to exist.
//
// It will check immediately and then every interval duration until:
// - The file exists
// - The context is canceled
// - The maximum number of tries is reached
//
// Parameters:
//   - ctx: Context for cancellation
//   - interval: Time to wait between tries
//   - maxTries: Maximum number of times to try (including the immediate try)
//   - filePath: Path to the file to check
//
// Returns:
//   - error: nil if the file exists, otherwise an error explaining why it failed
//
// Example:
//
//	// Wait for a log file to be created, checking every second, up to 10 tries
//	err := util.WaitForFile(ctx, time.Second, 10, "/var/log/app.log")
func WaitForFile(ctx context.Context, interval time.Duration, maxTries uint, filePath string) error {
	return WaitFor(ctx, interval, maxTries, func() bool {
		_, err := os.Stat(filePath)
		return err == nil
	})
}
