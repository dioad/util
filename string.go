package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
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

// SensitiveString Not 'secure' still uses a string as a base type
// however does protect against accidental exposure in logs
type MaskedString struct {
	string
	Config MaskedConfig
}

func MaskedStringDecodeHook(from, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() != reflect.String || to != reflect.TypeOf(MaskedString{}) {
		return data, nil
	}

	return NewMaskedString(data.(string)), nil
}

// type U struct {
// 	Type string
// 	Name string
// }
//
// type B struct {
// 	u U `json:"u"`
// }

type MaskedConfig struct {
	PrefixCount      uint
	SuffixCount      uint
	Mask             string
	MinMask          uint
	ObfuscateLength  bool
	ObfuscatedLength uint
}

func (s *MaskedString) String() string {
	l := uint(len(s.string))
	if s.Config.ObfuscateLength {
		l = s.Config.ObfuscatedLength
	}

	prefixCount := s.Config.PrefixCount
	if prefixCount > l {
		prefixCount = 0
	}

	suffixCount := s.Config.SuffixCount
	if suffixCount > l {
		suffixCount = 0
	}

	unmaskedCharCount := prefixCount + suffixCount

	charsToMask := l - unmaskedCharCount

	minMask := s.Config.MinMask

	if minMask != 0 && minMask > charsToMask {
		prefixCount = 0
		suffixCount = 0
	}

	if unmaskedCharCount >= uint(len(s.string)) {
		prefixCount = 0
		suffixCount = 0
	}

	prefix := ""
	if prefixCount > 0 {
		prefix = s.string[:prefixCount]
	}

	suffix := ""
	if suffixCount > 0 {
		leadingChars := len(s.string) - int(suffixCount)
		suffix = s.string[leadingChars:]
	}

	paddingCount := l - (prefixCount + suffixCount)

	maskChar := "*"
	if s.Config.Mask != "" {
		maskChar = s.Config.Mask
	}

	mask := strings.Repeat(maskChar, int(paddingCount))

	return fmt.Sprintf("%s%s%s", prefix, mask, suffix)
}

func (s *MaskedString) MaskedString() string {
	return s.string
}

func (s *MaskedString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	s.string = str
	return nil
}

// NewMaskedString creates a new masked string
func NewMaskedString(s string) *MaskedString {
	baseLength := int(1.5 * float32(len(s)))
	randomLength := rand.Intn(baseLength)

	m := &MaskedString{
		string: s,
	}
	m.Config.ObfuscatedLength = uint(randomLength)

	return m
}
