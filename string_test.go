package util

import (
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

func TestSensitiveString(t *testing.T) {
	s := NewSensitiveString("test")
	if s.String() != "********" {
		t.Errorf("expected '********' got '%s'", s.String())
	}
	if fmt.Sprintf("%v", s) != "********" {
		t.Errorf("expected '********' got '%s'", s)
	}
	if s.SensitiveString() != "test" {
		t.Errorf("expected 'test' got '%s'", s.SensitiveString())
	}
}
