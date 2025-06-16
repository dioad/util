package util

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"text/template"
)

// ExpandStringTemplate expands a string template with data.
func ExpandStringTemplate(templateString string, data any) (string, error) {
	tmpl, err := template.New("tmpl").Parse(templateString)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// MaskedString provides a way to handle sensitive string values while preventing
// accidental exposure in logs or user interfaces. When the string is displayed
// or logged, it shows a masked version instead of the actual sensitive value.
// Note: This is not a secure storage mechanism as it still uses a string as the
// underlying type, but it helps prevent accidental exposure.
type MaskedString struct {
	string
	Config MaskedConfig
}

// MaskedStringDecodeHook is a Viper decode hook that converts a string to a MaskedString.
// Use this with viper.DecodeHook() when unmarshaling configuration.
//
// Example usage with Viper:
//
//	viper.Unmarshal(&config, viper.DecodeHook(util.MaskedStringDecodeHook))
func MaskedStringDecodeHook(from, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() != reflect.String || to != reflect.TypeOf(MaskedString{}) {
		return data, nil
	}

	return NewMaskedString(data.(string)), nil
}

// MaskedConfig defines how a string should be masked when displayed.
type MaskedConfig struct {
	// PrefixCount is the number of characters to show at the beginning of the string
	PrefixCount uint
	// SuffixCount is the number of characters to show at the end of the string
	SuffixCount uint
	// Mask is the character to use for masking (defaults to "*")
	Mask string
	// MinMask is the minimum number of mask characters to show
	// If the number of characters to mask is less than MinMask,
	// the entire string will be masked
	MinMask uint
	// ObfuscateLength when true will use ObfuscatedLength instead of the actual length
	ObfuscateLength bool
	// ObfuscatedLength is the length to use when ObfuscateLength is true
	ObfuscatedLength uint
}

// String returns a masked version of the string according to the configuration.
// This method is called when the string is printed or logged.
func (s *MaskedString) String() string {
	// Determine the effective length to use
	effectiveLength := uint(len(s.string))
	if s.Config.ObfuscateLength {
		effectiveLength = s.Config.ObfuscatedLength
	}

	// Validate prefix and suffix counts
	prefixCount := s.Config.PrefixCount
	if prefixCount > effectiveLength {
		prefixCount = 0
	}

	suffixCount := s.Config.SuffixCount
	if suffixCount > effectiveLength {
		suffixCount = 0
	}

	// Check if we need to mask the entire string
	totalUnmasked := prefixCount + suffixCount
	charsToMask := effectiveLength - totalUnmasked

	// Apply minimum mask requirement if needed
	if s.Config.MinMask > 0 && s.Config.MinMask > charsToMask {
		prefixCount = 0
		suffixCount = 0
	}

	// If we're trying to show more characters than exist, mask everything
	if totalUnmasked >= uint(len(s.string)) {
		prefixCount = 0
		suffixCount = 0
	}

	// Extract prefix if needed
	prefix := ""
	if prefixCount > 0 {
		prefix = s.string[:prefixCount]
	}

	// Extract suffix if needed
	suffix := ""
	if suffixCount > 0 {
		leadingChars := len(s.string) - int(suffixCount)
		suffix = s.string[leadingChars:]
	}

	// Calculate how many mask characters to show
	paddingCount := effectiveLength - (prefixCount + suffixCount)

	// Determine which mask character to use
	maskChar := "*"
	if s.Config.Mask != "" {
		maskChar = s.Config.Mask
	}

	// Create the mask string
	mask := strings.Repeat(maskChar, int(paddingCount))

	// Combine prefix, mask, and suffix
	return fmt.Sprintf("%s%s%s", prefix, mask, suffix)
}

// UnmaskedString returns the original, unmasked string.
// Use this method with caution as it exposes the sensitive value.
func (s *MaskedString) UnmaskedString() string {
	return s.string
}

// MarshalJSON implements the json.Marshaler interface.
// It marshals the unmasked string value.
func (s *MaskedString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.string)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// It unmarshals a JSON string into a MaskedString.
func (s *MaskedString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	s.string = str
	return nil
}

func maskedStringLength(s string) uint {
	baseLength := int(1.5 * float32(len(s)))
	if baseLength == 0 {
		baseLength = 8
	}

	// Use crypto/rand for more secure random number generation
	n, err := rand.Int(rand.Reader, big.NewInt(int64(baseLength)))
	if err != nil {
		// Fallback to a simple calculation if random generation fails
		return uint(baseLength / 2)
	}
	return uint(n.Int64())
}

// NewMaskedString creates a new MaskedString with default configuration.
// By default, the string will be completely masked with asterisks (*) and
// will have a randomized length to prevent length-based information leakage.
//
// Example usage:
//
//	password := util.NewMaskedString("sensitive-password")
//	fmt.Println(password) // Prints something like "**********"
//
//	// Configure to show first and last character
//	password.Config.PrefixCount = 1
//	password.Config.SuffixCount = 1
//	fmt.Println(password) // Prints something like "s********d"
//
//	// Use a different mask character
//	password.Config.Mask = "•"
//	fmt.Println(password) // Prints something like "s••••••••d"
//
//	// Access the original value when needed (use with caution)
//	originalPassword := password.UnmaskedString()
func NewMaskedString(s string) *MaskedString {
	m := &MaskedString{
		string: s,
		Config: MaskedConfig{
			ObfuscateLength:  true,
			ObfuscatedLength: maskedStringLength(s),
		},
	}

	return m
}
