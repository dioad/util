package util

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

// LookupEnvWithDefault returns the value of the environment variable named by the key.
// If the variable is not present, it returns the defaultValue.
//
// Example:
//
//	// Get database host from environment or use localhost as default
//	dbHost := util.LookupEnvWithDefault("DB_HOST", "localhost")
func LookupEnvWithDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

// LookupEnvBool returns the boolean value of the environment variable named by the key.
// It returns an error if the variable is not present or cannot be parsed as a boolean.
//
// Valid boolean values are: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
//
// Example:
//
//	// Check if debug mode is enabled
//	debug, err := util.LookupEnvBool("DEBUG_MODE")
//	if err != nil {
//	    // Handle error or use default
//	    debug = false
//	}
func LookupEnvBool(key string) (bool, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return false, fmt.Errorf("environment variable %s is not set", key)
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("environment variable %s is not a valid boolean: %w", key, err)
	}

	return b, nil
}

// LookupEnvURL returns the URL value of the environment variable named by the key.
// It returns nil, nil if the variable is not present.
// It returns nil, error if the variable cannot be parsed as a URL.
//
// Example:
//
//	// Get API endpoint URL from environment
//	apiURL, err := util.LookupEnvURL("API_ENDPOINT")
//	if err != nil {
//	    return fmt.Errorf("invalid API endpoint URL: %w", err)
//	}
//	if apiURL == nil {
//	    // Use default URL if not set
//	    apiURL, _ = url.Parse("https://api.example.com")
//	}
func LookupEnvURL(key string) (*url.URL, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return nil, nil // Variable not set, not an error
	}

	parsedURL, err := url.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("unable to parse environment variable %s as URL: %w", key, err)
	}

	return parsedURL, nil
}

// LookupEnvInt returns the integer value of the environment variable named by the key.
// It returns an error if the variable is not present or cannot be parsed as an integer.
//
// Example:
//
//	// Get port number from environment
//	port, err := util.LookupEnvInt("PORT")
//	if err != nil {
//	    // Handle error or use default
//	    port = 8080
//	}
func LookupEnvInt(key string) (int, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, fmt.Errorf("environment variable %s is not set", key)
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not a valid integer: %w", key, err)
	}

	return i, nil
}

// LookupEnvFloat64 returns the float64 value of the environment variable named by the key.
// It returns an error if the variable is not present or cannot be parsed as a float64.
//
// Example:
//
//	// Get a threshold value from environment
//	threshold, err := util.LookupEnvFloat64("THRESHOLD")
//	if err != nil {
//	    // Handle error or use default
//	    threshold = 0.5
//	}
func LookupEnvFloat64(key string) (float64, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, fmt.Errorf("environment variable %s is not set", key)
	}

	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not a valid float: %w", key, err)
	}

	return f, nil
}

// LookupEnvDuration returns the duration value of the environment variable named by the key.
// It returns an error if the variable is not present or cannot be parsed as a duration.
//
// Example:
//
//	// Get timeout duration from environment
//	timeout, err := util.LookupEnvDuration("TIMEOUT")
//	if err != nil {
//	    // Handle error or use default
//	    timeout = 30 * time.Second
//	}
func LookupEnvDuration(key string) (time.Duration, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, fmt.Errorf("environment variable %s is not set", key)
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not a valid duration: %w", key, err)
	}

	return d, nil
}
