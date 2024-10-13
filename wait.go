package util

import (
	"fmt"
	"time"
)

// WaitFor waits for a function to return true, it will check every interval seconds up until max seconds.
func WaitFor(interval time.Duration, maxTries uint, op func() bool) error {
	var i uint
	for i = 0; i < maxTries; i++ {
		if op() {
			return nil
		}
		time.Sleep(interval)
	}
	return fmt.Errorf("condition not met")
}

// WaitForNilError waits for a function to return a nil error, it will check every interval seconds up until max seconds.
func WaitForNilError(interval time.Duration, maxTries uint, op func() error) error {
	return WaitFor(interval, maxTries, func() bool {
		return op() == nil
	})
}

// WaitForReturn waits for a function to return a non-nil value, it will check every interval seconds up until max seconds.
// The function returns the value and error returned by the function.
// If maxTries is 0, it will only try once (it will set maxTries internally to 1).
func WaitForReturn[T any](interval time.Duration, maxTries uint, op func() (*T, error)) (*T, error) {
	var i uint

	if maxTries == 0 {
		maxTries = 1
	}

	for i = 0; i < maxTries; i++ {
		resp, err := op()
		if err == nil {
			return resp, nil
		}
		time.Sleep(interval)
	}
	return nil, fmt.Errorf("condition not met")
}
