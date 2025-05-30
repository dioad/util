package util

import (
	"net/url"
	"testing"
)

func mockLookupEnv(lookupKey, result string) envLookup {
	return func(key string) (string, bool) {
		if key != lookupKey {
			return "", false
		}
		return result, true
	}
}

func TestLookupEnvWithDefault(t *testing.T) {
	tests := []struct {
		key          string
		defaultValue string
		lookupFunc   envLookup
		expected     string
	}{
		{
			key:          "TEST_KEY",
			lookupFunc:   mockLookupEnv("TEST_KEY", "value"),
			defaultValue: "defaultValue",
			expected:     "value",
		},
		{
			key:          "TEST_KEY_NO_VALUE",
			lookupFunc:   mockLookupEnv("TEST_KEY", "value"),
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
		lookupFunc envLookup
		expected   bool
	}{
		{
			key:        "TEST_KEY",
			lookupFunc: mockLookupEnv("TEST_KEY", "true"),
			expected:   true,
		},
		{
			key:        "TEST_KEY",
			lookupFunc: mockLookupEnv("TEST_KEY", "TRUE"),
			expected:   true,
		},
		{
			key:        "TEST_KEY",
			lookupFunc: mockLookupEnv("TEST_KEY", "1"),
			expected:   true,
		},
		{
			key:        "TEST_KEY",
			lookupFunc: mockLookupEnv("TEST_NO_KEY", "asdf"),
			expected:   false,
		},
		{
			key:        "TEST_KEY",
			lookupFunc: mockLookupEnv("TEST_KEY", "asdf"),
			expected:   false,
		},
	}

	for _, test := range tests {
		if value, err := lookupEnvBool(test.lookupFunc, test.key); value != test.expected {
			if err != nil && test.expected {
				t.Fatalf("failed to lookup %v, got %v", test.expected, err)
			} else if err == nil && !test.expected {
				t.Fatalf("expected error for key %s, got nil", test.key)
			} else if err != nil && !test.expected {
				continue // This is expected for invalid boolean values
			}
			t.Fatalf("expected %v, got %v", test.expected, value)
		}
	}
}

func MustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func TestLookupEnvURL(t *testing.T) {
	tests := []struct {
		key           string
		lookupFunc    envLookup
		expectedValue *url.URL
		errorExpected bool
	}{
		{
			key:           "TEST_KEY",
			lookupFunc:    mockLookupEnv("TEST_KEY", "https://asdf/asdf"),
			expectedValue: MustParseURL("https://asdf/asdf"),
			errorExpected: false,
		},
		{
			key:           "TEST_KEY_INVALID_VALUE",
			lookupFunc:    mockLookupEnv("TEST_KEY_INVALID_VALUE", "asdf\nasdf"),
			expectedValue: nil,
			errorExpected: true,
		},
		{
			key:           "TEST_KEY_NO_VALUE",
			lookupFunc:    mockLookupEnv("TEST_KEY", "https://asdf/asdf"),
			expectedValue: nil,
			errorExpected: false,
		},
	}

	for _, test := range tests {
		value, err := lookupEnvURL(test.lookupFunc, test.key)

		if err != nil && !test.errorExpected {
			t.Fatalf("failed to lookup %v, got %v", test.expectedValue, err)
		}

		if err == nil && test.errorExpected {
			t.Fatalf("expected error, got %v", value)
		}

		if value == nil && test.expectedValue != nil {
			t.Fatalf("expected %v, got nil", test.expectedValue)
		}

		if value != nil && test.expectedValue == nil {
			t.Fatalf("expected nil, got %v", value)
		}

		if value != nil && test.expectedValue != nil {
			if value.String() != test.expectedValue.String() {
				t.Fatalf("expected %v, got %v", test.expectedValue, value)
			}
		}
	}
}
