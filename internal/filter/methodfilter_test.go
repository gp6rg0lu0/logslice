package filter

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makeMethodLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]string{field: value})
}

func TestNewMethodFilter_EmptyField(t *testing.T) {
	_, err := NewMethodFilter("", []string{"GET"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewMethodFilter_NoMethods(t *testing.T) {
	_, err := NewMethodFilter("method", []string{})
	if err == nil {
		t.Fatal("expected error for empty methods slice")
	}
}

func TestNewMethodFilter_BlankMethod(t *testing.T) {
	_, err := NewMethodFilter("method", []string{"GET", "  "})
	if err == nil {
		t.Fatal("expected error for blank method entry")
	}
}

func TestMethodFilter_Accessors(t *testing.T) {
	f, err := NewMethodFilter("method", []string{"GET", "POST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "method" {
		t.Errorf("Field() = %q, want %q", f.Field(), "method")
	}
	if len(f.Methods()) != 2 {
		t.Errorf("Methods() len = %d, want 2", len(f.Methods()))
	}
}

func TestMethodFilter_NilLine(t *testing.T) {
	f, _ := NewMethodFilter("method", []string{"GET"})
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestMethodFilter_MissingField(t *testing.T) {
	f, _ := NewMethodFilter("method", []string{"GET"})
	line := makeMethodLogLine("other", "GET")
	if f.Match(line) {
		t.Error("expected false when field is absent")
	}
}

func TestMethodFilter_MatchesExact(t *testing.T) {
	f, _ := NewMethodFilter("method", []string{"GET", "POST"})
	for _, m := range []string{"GET", "POST"} {
		line := makeMethodLogLine("method", m)
		if !f.Match(line) {
			t.Errorf("expected match for method %q", m)
		}
	}
}

func TestMethodFilter_CaseInsensitive(t *testing.T) {
	f, _ := NewMethodFilter("method", []string{"DELETE"})
	for _, m := range []string{"delete", "Delete", "DELETE"} {
		line := makeMethodLogLine("method", m)
		if !f.Match(line) {
			t.Errorf("expected case-insensitive match for %q", m)
		}
	}
}

func TestMethodFilter_NoMatch(t *testing.T) {
	f, _ := NewMethodFilter("method", []string{"GET", "POST"})
	line := makeMethodLogLine("method", "DELETE")
	if f.Match(line) {
		t.Error("expected no match for DELETE when filter is GET/POST")
	}
}
