package util

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

type envLookup func(string) (string, bool)

// lookupEnvWithDefault is a helper function that returns a value from an environment variable with a default value
func lookupEnvWithDefault(lookup envLookup, key, defaultValue string) string {
	if value, ok := lookup(key); ok {
		return value
	}
	return defaultValue
}

// lookupEnvBool is a helper function that returns a boolean value from an environment variable
func lookupEnvBool(lookup envLookup, key string) bool {
	if value, ok := lookup(key); ok {
		return strings.ToLower(value) == "true"
	}
	return false
}

// lookupEnvURL is a helper function that returns a URL from an environment variable
func lookupEnvURL(lookup envLookup, key string) (*url.URL, error) {
	if value, ok := lookup(key); ok {
		tmpURL, err := url.Parse(value)
		if err != nil {
			return nil, fmt.Errorf("unable to parse %v as URL: %w", value, err)
		}
		return tmpURL, nil
	}
	return nil, nil
}

// LookupEnvWithDefault is a wrapper around os.LookupEnv that returns a default value if the environment variable is not set
func LookupEnvWithDefault(key, defaultValue string) string {
	return lookupEnvWithDefault(os.LookupEnv, key, defaultValue)
}

// LookupEnvBool is a wrapper around os.LookupEnv that returns a boolean value
func LookupEnvBool(key string) bool {
	return lookupEnvBool(os.LookupEnv, key)
}

// LookupEnvURL is a wrapper around os.LookupEnv that returns a URL
func LookupEnvURL(key string) (*url.URL, error) {
	return lookupEnvURL(os.LookupEnv, key)
}
