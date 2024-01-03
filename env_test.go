package util

import "testing"

func mockLookupEnv(lookupKey, result string, exists bool) EnvLookup {
	return func(key string) (string, bool) {
		if key != lookupKey {
			return "", false
		}
		return result, exists
	}
}

func TestLookupEnvWithDefault(t *testing.T) {
	tests := []struct {
		key          string
		defaultValue string
		lookupFunc   EnvLookup
		expected     string
	}{
		{
			key:          "TEST_KEY",
			lookupFunc:   mockLookupEnv("TEST_KEY", "value", true),
			defaultValue: "defaultValue",
			expected:     "value",
		},
		{
			key:          "TEST_KEY",
			lookupFunc:   mockLookupEnv("TEST_KEY", "value", false),
			defaultValue: "defaultValue",
			expected:     "defaultValue",
		},
	}

	for _, test := range tests {
		if value := lookupEnvWithDefault(test.lookupFunc, test.key, test.defaultValue); value != test.expected {
			t.Fatalf("expected %v, got %v", test.expected, value)
		}
	}
}

func TestLookupEnvBool(t *testing.T) {
	tests := []struct {
		key        string
		lookupFunc EnvLookup
		expected   bool
	}{
		{
			key:        "TEST_KEY",
			lookupFunc: mockLookupEnv("TEST_KEY", "true", true),
			expected:   true,
		},
		{
			key:        "TEST_KEY",
			lookupFunc: mockLookupEnv("TEST_KEY", "asdf", false),
			expected:   false,
		},
		{
			key:        "TEST_KEY",
			lookupFunc: mockLookupEnv("TEST_KEY", "asdf", true),
			expected:   false,
		},
	}

	for _, test := range tests {
		if value := lookupEnvBool(test.lookupFunc, test.key); value != test.expected {
			t.Fatalf("expected %v, got %v", test.expected, value)
		}
	}
}
