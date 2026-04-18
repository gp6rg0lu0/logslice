package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeMultiValueLogLine(field, value string) *parser.LogLine {
	fields := map[string]interface{}{field: value}
	return parser.NewLogLine(fields)
}

func TestNewMultiValueFilter_EmptyField(t *testing.T) {
	_, err := NewMultiValueFilter("", []string{"a"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewMultiValueFilter_NoValues(t *testing.T) {
	_, err := NewMultiValueFilter("level", []string{})
	if err == nil {
		t.Fatal("expected error for empty values")
	}
}

func TestNewMultiValueFilter_Valid(t *testing.T) {
	f, err := NewMultiValueFilter("level", []string{"info", "warn"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "level" {
		t.Errorf("expected field 'level', got %q", f.Field())
	}
	if len(f.Values()) != 2 {
		t.Errorf("expected 2 values, got %d", len(f.Values()))
	}
}

func TestMultiValueFilter_NilLine(t *testing.T) {
	f, _ := NewMultiValueFilter("level", []string{"info"})
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestMultiValueFilter_MissingField(t *testing.T) {
	f, _ := NewMultiValueFilter("level", []string{"info"})
	line := makeMultiValueLogLine("other", "info")
	if f.Match(line) {
		t.Error("expected false for missing field")
	}
}

func TestMultiValueFilter_MatchesAllowedValue(t *testing.T) {
	f, _ := NewMultiValueFilter("level", []string{"info", "warn", "error"})
	for _, v := range []string{"info", "warn", "error"} {
		line := makeMultiValueLogLine("level", v)
		if !f.Match(line) {
			t.Errorf("expected match for value %q", v)
		}
	}
}

func TestMultiValueFilter_RejectsOtherValue(t *testing.T) {
	f, _ := NewMultiValueFilter("level", []string{"info", "warn"})
	line := makeMultiValueLogLine("level", "debug")
	if f.Match(line) {
		t.Error("expected no match for 'debug'")
	}
}

func TestMultiValueFilter_TrimsSpaces(t *testing.T) {
	f, _ := NewMultiValueFilter("level", []string{" info ", "warn"})
	line := makeMultiValueLogLine("level", "info")
	if !f.Match(line) {
		t.Error("expected match after trimming spaces in allowed values")
	}
}
