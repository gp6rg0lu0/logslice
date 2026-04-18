package filter

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makeNullLogLine(fields map[string]string) *parser.LogLine {
	return parser.NewLogLine(fields)
}

func TestNewNullFilter_EmptyField(t *testing.T) {
	_, err := NewNullFilter("", true)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewNullFilter_Valid(t *testing.T) {
	f, err := NewNullFilter("error", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "error" {
		t.Errorf("expected field 'error', got %q", f.Field())
	}
	if !f.MustBeNull() {
		t.Error("expected MustBeNull to be true")
	}
}

func TestNullFilter_NilLine(t *testing.T) {
	f, _ := NewNullFilter("msg", true)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestNullFilter_FieldMissing_MustBeNull(t *testing.T) {
	f, _ := NewNullFilter("missing", true)
	line := makeNullLogLine(map[string]string{"other": "value"})
	if !f.Match(line) {
		t.Error("expected match when field is missing and mustBeNull=true")
	}
}

func TestNullFilter_FieldPresent_MustBeNull(t *testing.T) {
	f, _ := NewNullFilter("msg", true)
	line := makeNullLogLine(map[string]string{"msg": "hello"})
	if f.Match(line) {
		t.Error("expected no match when field is present and mustBeNull=true")
	}
}

func TestNullFilter_FieldMissing_MustNotBeNull(t *testing.T) {
	f, _ := NewNullFilter("msg", false)
	line := makeNullLogLine(map[string]string{"other": "x"})
	if f.Match(line) {
		t.Error("expected no match when field is missing and mustBeNull=false")
	}
}

func TestNullFilter_FieldPresent_MustNotBeNull(t *testing.T) {
	f, _ := NewNullFilter("msg", false)
	line := makeNullLogLine(map[string]string{"msg": "hello"})
	if !f.Match(line) {
		t.Error("expected match when field is present and mustBeNull=false")
	}
}

func TestNullFilter_EmptyValue_TreatedAsNull(t *testing.T) {
	f, _ := NewNullFilter("msg", true)
	line := makeNullLogLine(map[string]string{"msg": ""})
	if !f.Match(line) {
		t.Error("expected match when field is empty string and mustBeNull=true")
	}
}
