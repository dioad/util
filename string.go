package util

import (
	"bytes"
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
type SensitiveString struct {
	string
}

func (s SensitiveString) String() string {
	return "********"
}

func (s SensitiveString) SensitiveString() string {
	return s.string
}

func NewSensitiveString(s string) *SensitiveString {
	return &SensitiveString{s}
}
