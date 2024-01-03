package util

import "os"

type EnvLookup func(string) (string, bool)

func lookupEnvWithDefault(lookup EnvLookup, key, defaultValue string) string {
	if value, ok := lookup(key); ok {
		return value
	}
	return defaultValue
}

func lookupEnvBool(lookup EnvLookup, key string) bool {
	if value, ok := lookup(key); ok {
		if value == "true" {
			return true
		}
		return false
	}
	return false
}

func LookupEnvWithDefault(key, defaultValue string) string {
	return lookupEnvWithDefault(os.LookupEnv, key, defaultValue)
}

func LookupEnvBool(key string) bool {
	return lookupEnvBool(os.LookupEnv, key)
}
