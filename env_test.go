package util

import (
	"net/url"
	"os"
	"testing"
	"time"
)

func TestLookupEnvWithDefault(t *testing.T) {
	t.Run("existing variable", func(t *testing.T) {
		os.Setenv("TEST_KEY", "value")
		defer os.Unsetenv("TEST_KEY")
		if value := LookupEnvWithDefault("TEST_KEY", "defaultValue"); value != "value" {
			t.Fatalf("expected value, got %v", value)
		}
	})

	t.Run("non-existing variable", func(t *testing.T) {
		os.Unsetenv("TEST_KEY")
		if value := LookupEnvWithDefault("TEST_KEY", "defaultValue"); value != "defaultValue" {
			t.Fatalf("expected defaultValue, got %v", value)
		}
	})
}

func TestLookupEnvBool(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		set      bool
		expected bool
		wantErr  bool
	}{
		{
			name:     "valid true",
			key:      "TEST_KEY_TRUE",
			value:    "true",
			set:      true,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "valid TRUE",
			key:      "TEST_KEY_TRUE_UPPER",
			value:    "TRUE",
			set:      true,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "valid 1",
			key:      "TEST_KEY_1",
			value:    "1",
			set:      true,
			expected: true,
			wantErr:  false,
		},
		{
			name:    "not set",
			key:     "TEST_KEY_UNSET",
			set:     false,
			wantErr: true,
		},
		{
			name:    "invalid value",
			key:     "TEST_KEY_INVALID",
			value:   "not-a-bool",
			set:     true,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			value, err := LookupEnvBool(tt.key)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error expectation mismatch: wantErr %v, got %v", tt.wantErr, err)
			}
			if !tt.wantErr && value != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, value)
			}
		})
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
		name          string
		key           string
		value         string
		set           bool
		expectedValue *url.URL
		wantErr       bool
	}{
		{
			name:          "valid URL",
			key:           "TEST_KEY_URL",
			value:         "https://example.com/path",
			set:           true,
			expectedValue: MustParseURL("https://example.com/path"),
			wantErr:       false,
		},
		{
			name:    "invalid URL",
			key:     "TEST_KEY_INVALID_URL",
			value:   "asdf\nasdf",
			set:     true,
			wantErr: true,
		},
		{
			name:          "not set",
			key:           "TEST_KEY_UNSET_URL",
			set:           false,
			expectedValue: nil,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			value, err := LookupEnvURL(tt.key)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error expectation mismatch: wantErr %v, got %v", tt.wantErr, err)
			}
			if !tt.wantErr {
				if value == nil && tt.expectedValue != nil {
					t.Fatalf("expected %v, got nil", tt.expectedValue)
				}
				if value != nil && tt.expectedValue == nil {
					t.Fatalf("expected nil, got %v", value)
				}
				if value != nil && tt.expectedValue != nil && value.String() != tt.expectedValue.String() {
					t.Fatalf("expected %v, got %v", tt.expectedValue, value)
				}
			}
		})
	}
}

func TestLookupEnvInt(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		value         string
		set           bool
		expectedValue int
		wantErr       bool
	}{
		{
			name:          "valid integer",
			key:           "TEST_KEY_INT",
			value:         "42",
			set:           true,
			expectedValue: 42,
			wantErr:       false,
		},
		{
			name:    "invalid integer",
			key:     "TEST_KEY_INVALID_INT",
			value:   "not-an-integer",
			set:     true,
			wantErr: true,
		},
		{
			name:    "not set",
			key:     "TEST_KEY_UNSET_INT",
			set:     false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			value, err := LookupEnvInt(tt.key)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error expectation mismatch: wantErr %v, got %v", tt.wantErr, err)
			}
			if !tt.wantErr && value != tt.expectedValue {
				t.Fatalf("expected %d, got %d", tt.expectedValue, value)
			}
		})
	}
}

func TestLookupEnvFloat64(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		value         string
		set           bool
		expectedValue float64
		wantErr       bool
	}{
		{
			name:          "valid float",
			key:           "TEST_KEY_FLOAT",
			value:         "3.14159",
			set:           true,
			expectedValue: 3.14159,
			wantErr:       false,
		},
		{
			name:    "invalid float",
			key:     "TEST_KEY_INVALID_FLOAT",
			value:   "not-a-float",
			set:     true,
			wantErr: true,
		},
		{
			name:    "not set",
			key:     "TEST_KEY_UNSET_FLOAT",
			set:     false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			value, err := LookupEnvFloat64(tt.key)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error expectation mismatch: wantErr %v, got %v", tt.wantErr, err)
			}
			if !tt.wantErr && value != tt.expectedValue {
				t.Fatalf("expected %f, got %f", tt.expectedValue, value)
			}
		})
	}
}

func TestLookupEnvDuration(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		value         string
		set           bool
		expectedValue time.Duration
		wantErr       bool
	}{
		{
			name:          "valid duration",
			key:           "TEST_KEY_DURATION",
			value:         "5s",
			set:           true,
			expectedValue: 5 * time.Second,
			wantErr:       false,
		},
		{
			name:    "invalid duration",
			key:     "TEST_KEY_INVALID_DURATION",
			value:   "not-a-duration",
			set:     true,
			wantErr: true,
		},
		{
			name:    "not set",
			key:     "TEST_KEY_UNSET_DURATION",
			set:     false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			value, err := LookupEnvDuration(tt.key)
			if (err != nil) != tt.wantErr {
				t.Fatalf("error expectation mismatch: wantErr %v, got %v", tt.wantErr, err)
			}
			if !tt.wantErr && value != tt.expectedValue {
				t.Fatalf("expected %v, got %v", tt.expectedValue, value)
			}
		})
	}
}
