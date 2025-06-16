package util

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestExpandStringTemplate(t *testing.T) {
	type testStruct struct {
		One string
		Two string
	}
	data := testStruct{
		One: "one",
		Two: "two",
	}
	templateString := "{{.One}} {{.Two}}"
	result, err := ExpandStringTemplate(templateString, data)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if result != "one two" {
		t.Errorf("expected 'one two' got '%s'", result)
	}
}

func TestMaskedStringJSON(t *testing.T) {
	// Test JSON marshaling and unmarshaling
	original := "sensitive-data"
	ms := NewMaskedString(original)

	// Test MarshalJSON
	jsonData, err := json.Marshal(ms)
	if err != nil {
		t.Errorf("unexpected error marshaling MaskedString: %v", err)
	}

	// The JSON should be the original string in quotes
	expected := fmt.Sprintf("\"%s\"", original)
	if string(jsonData) != expected {
		t.Errorf("expected JSON %s, got %s", expected, string(jsonData))
	}

	// Test UnmarshalJSON
	var newMS MaskedString
	err = json.Unmarshal(jsonData, &newMS)
	if err != nil {
		t.Errorf("unexpected error unmarshaling MaskedString: %v", err)
	}

	if newMS.UnmaskedString() != original {
		t.Errorf("expected unmarshaled value %s, got %s", original, newMS.UnmaskedString())
	}
}

func TestMaskedString(t *testing.T) {
	tests := []struct {
		name     string
		cfg      MaskedConfig
		str      string
		expected string
	}{
		{
			name:     "empty",
			cfg:      MaskedConfig{},
			str:      "test",
			expected: "****",
		},
		{
			name: "custom mask",
			cfg: MaskedConfig{
				Mask: "X",
			},
			str:      "test",
			expected: "XXXX",
		},
		{
			name: "prefix",
			cfg: MaskedConfig{
				PrefixCount: 1,
			},
			str:      "test",
			expected: "t***",
		},
		{
			name: "suffix",
			cfg: MaskedConfig{
				SuffixCount: 1,
			},
			str:      "test",
			expected: "***t",
		},
		{
			name: "prefix and suffix",
			cfg: MaskedConfig{
				PrefixCount: 1,
				SuffixCount: 1,
			},
			str:      "test",
			expected: "t**t",
		},
		{
			name: "prefix and suffix and mask",
			cfg: MaskedConfig{
				PrefixCount: 1,
				SuffixCount: 1,
				Mask:        "X",
			},
			str:      "test",
			expected: "tXXt",
		},
		{
			name: "prefix and suffix and mask",
			cfg: MaskedConfig{
				PrefixCount: 1,
				SuffixCount: 2,
				Mask:        "X",
			},
			str:      "test",
			expected: "tXst",
		},
		{
			name: "prefix and suffix and mask",
			cfg: MaskedConfig{
				PrefixCount: 2,
				SuffixCount: 2,
				Mask:        "X",
			},
			str:      "test",
			expected: "XXXX",
		},
		{
			name: "prefix and suffix and mask",
			cfg: MaskedConfig{
				PrefixCount: 5,
				SuffixCount: 5,
				Mask:        "X",
			},
			str:      "test",
			expected: "XXXX",
		},
		{
			name: "prefix and suffix and mask",
			cfg: MaskedConfig{
				PrefixCount: 1,
				SuffixCount: 2,
				MinMask:     2,
				Mask:        "X",
			},
			str:      "test",
			expected: "XXXX",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMaskedString(tt.str)
			s.Config = tt.cfg
			if s.String() != tt.expected {
				t.Errorf("expected '%s' got '%s'", tt.expected, s.String())
			}
			if fmt.Sprintf("%v", s) != tt.expected {
				t.Errorf("expected '%s' got '%s'", tt.expected, s)
			}
			if s.UnmaskedString() != tt.str {
				t.Errorf("expected '%s' got '%s'", tt.str, s.UnmaskedString())
			}
		})
	}
}

func TestMaskedStringWithObfuscatedLength(t *testing.T) {
	tests := []struct {
		name     string
		cfg      MaskedConfig
		str      string
		expected string
	}{
		{
			name: "less than string length",
			cfg: MaskedConfig{
				ObfuscateLength:  true,
				ObfuscatedLength: 3,
			},
			str:      "test",
			expected: "***",
		},
		{
			name: "equal to string length",
			cfg: MaskedConfig{
				ObfuscateLength:  true,
				ObfuscatedLength: 4,
			},
			str:      "test",
			expected: "****",
		},
		{
			name: "greater than string length",
			cfg: MaskedConfig{
				ObfuscateLength:  true,
				ObfuscatedLength: 8,
			},
			str:      "test",
			expected: "********",
		},
		{
			name: "greater than string length with PrefixCount",
			cfg: MaskedConfig{
				PrefixCount:      2,
				ObfuscateLength:  true,
				ObfuscatedLength: 8,
			},
			str:      "test",
			expected: "te******",
		},
		{
			name: "greater than string length with PrefixCount and SuffixCount",
			cfg: MaskedConfig{
				PrefixCount:      1,
				SuffixCount:      1,
				ObfuscateLength:  true,
				ObfuscatedLength: 8,
			},
			str:      "test",
			expected: "t******t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMaskedString(tt.str)
			s.Config = tt.cfg
			if s.String() != tt.expected {
				t.Errorf("expected '%s' got '%s'", tt.expected, s.String())
			}
			if fmt.Sprintf("%v", s) != tt.expected {
				t.Errorf("expected '%s' got '%s'", tt.expected, s)
			}
			if s.UnmaskedString() != tt.str {
				t.Errorf("expected '%s' got '%s'", tt.str, s.UnmaskedString())
			}
		})
	}
}
