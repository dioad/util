package util

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

// envLookup is a function type that looks up a key in an environment and returns its value and whether it exists.
// This abstraction allows for easier testing by mocking the environment.
type envLookup func(string) (string, bool)

// lookupEnvWithDefault returns a value from an environment variable or a default value if not set.
// It uses the provided lookup function to access the environment.
func lookupEnvWithDefault(lookup envLookup, key, defaultValue string) string {
	if value, ok := lookup(key); ok {
		return value
	}
	return defaultValue
}

// lookupEnvBool returns a boolean value from an environment variable.
// It uses the provided lookup function to access the environment.
// Returns an error if the environment variable is not set or cannot be parsed as a boolean.
func lookupEnvBool(lookup envLookup, key string) (bool, error) {
	value, ok := lookup(key)
	if !ok {
		return false, fmt.Errorf("environment variable %s is not set", key)
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("environment variable %s is not a valid boolean: %w", key, err)
	}

	return b, nil
}

// lookupEnvURL returns a URL from an environment variable.
// It uses the provided lookup function to access the environment.
// Returns nil, nil if the environment variable is not set.
// Returns nil, error if the environment variable cannot be parsed as a URL.
func lookupEnvURL(lookup envLookup, key string) (*url.URL, error) {
	value, ok := lookup(key)
	if !ok {
		return nil, nil // Variable not set, not an error
	}

	parsedURL, err := url.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("unable to parse environment variable %s value %q as URL: %w", key, value, err)
	}

	return parsedURL, nil
}

// LookupEnvWithDefault returns the value of the environment variable named by the key.
// If the variable is not present, it returns the defaultValue.
//
// Example:
//
//	// Get database host from environment or use localhost as default
//	dbHost := util.LookupEnvWithDefault("DB_HOST", "localhost")
func LookupEnvWithDefault(key, defaultValue string) string {
	return lookupEnvWithDefault(os.LookupEnv, key, defaultValue)
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
	return lookupEnvBool(os.LookupEnv, key)
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
//	    apiURL = MustParseURL("https://api.example.com")
//	}
func LookupEnvURL(key string) (*url.URL, error) {
	return lookupEnvURL(os.LookupEnv, key)
}

// lookupEnvInt returns an integer value from an environment variable.
// It uses the provided lookup function to access the environment.
// Returns an error if the environment variable is not set or cannot be parsed as an integer.
func lookupEnvInt(lookup envLookup, key string) (int, error) {
	value, ok := lookup(key)
	if !ok {
		return 0, fmt.Errorf("environment variable %s is not set", key)
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is not a valid integer: %w", key, err)
	}

	return i, nil
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
	return lookupEnvInt(os.LookupEnv, key)
}
