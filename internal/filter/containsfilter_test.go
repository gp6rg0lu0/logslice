package filter_test

import (
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeContainsLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]string{field: value})
}

func TestNewContainsFilter_EmptyField(t *testing.T) {
	_, err := filter.NewContainsFilter("", "foo", false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewContainsFilter_EmptySubstring(t *testing.T) {
	_, err := filter.NewContainsFilter("msg", "", false)
	if err == nil {
		t.Fatal("expected error for empty substring")
	}
}

func TestContainsFilter_Accessors(t *testing.T) {
	f, err := filter.NewContainsFilter("msg", "hello", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "msg" {
		t.Errorf("Field() = %q, want %q", f.Field(), "msg")
	}
	if f.Substring() != "hello" {
		t.Errorf("Substring() = %q, want %q", f.Substring(), "hello")
	}
	if !f.CaseFold() {
		t.Error("CaseFold() should be true")
	}
}

func TestContainsFilter_NilLine(t *testing.T) {
	f, _ := filter.NewContainsFilter("msg", "x", false)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestContainsFilter_MissingField(t *testing.T) {
	f, _ := filter.NewContainsFilter("msg", "hello", false)
	line := makeContainsLogLine("other", "hello world")
	if f.Match(line) {
		t.Error("expected false when field is missing")
	}
}

func TestContainsFilter_Match(t *testing.T) {
	f, _ := filter.NewContainsFilter("msg", "world", false)
	if !f.Match(makeContainsLogLine("msg", "hello world")) {
		t.Error("expected true when substring present")
	}
	if f.Match(makeContainsLogLine("msg", "hello there")) {
		t.Error("expected false when substring absent")
	}
}

func TestContainsFilter_CaseFold(t *testing.T) {
	f, _ := filter.NewContainsFilter("msg", "WORLD", true)
	if !f.Match(makeContainsLogLine("msg", "hello world")) {
		t.Error("expected case-insensitive match")
	}
	strict, _ := filter.NewContainsFilter("msg", "WORLD", false)
	if strict.Match(makeContainsLogLine("msg", "hello world")) {
		t.Error("expected case-sensitive mismatch")
	}
}
