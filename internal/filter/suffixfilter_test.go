package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeSuffixLogLine(data map[string]string) *parser.LogLine {
	return parser.NewLogLine(data)
}

func TestNewSuffixFilter_EmptyField(t *testing.T) {
	_, err := NewSuffixFilter("", ".log")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewSuffixFilter_EmptySuffix(t *testing.T) {
	_, err := NewSuffixFilter("file", "")
	if err == nil {
		t.Fatal("expected error for empty suffix")
	}
}

func TestSuffixFilter_Accessors(t *testing.T) {
	f, err := NewSuffixFilter("file", ".log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "file" {
		t.Errorf("expected field 'file', got %q", f.Field())
	}
	if f.Suffix() != ".log" {
		t.Errorf("expected suffix '.log', got %q", f.Suffix())
	}
}

func TestSuffixFilter_NilLine(t *testing.T) {
	f, _ := NewSuffixFilter("file", ".log")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestSuffixFilter_MissingField(t *testing.T) {
	f, _ := NewSuffixFilter("file", ".log")
	line := makeSuffixLogLine(map[string]string{"other": "app.log"})
	if f.Match(line) {
		t.Error("expected false when field is missing")
	}
}

func TestSuffixFilter_Matches(t *testing.T) {
	f, _ := NewSuffixFilter("file", ".log")
	line := makeSuffixLogLine(map[string]string{"file": "app.log"})
	if !f.Match(line) {
		t.Error("expected true for matching suffix")
	}
}

func TestSuffixFilter_NoMatch(t *testing.T) {
	f, _ := NewSuffixFilter("file", ".log")
	line := makeSuffixLogLine(map[string]string{"file": "app.json"})
	if f.Match(line) {
		t.Error("expected false for non-matching suffix")
	}
}

func TestSuffixFilter_ExactMatch(t *testing.T) {
	f, _ := NewSuffixFilter("msg", "error")
	line := makeSuffixLogLine(map[string]string{"msg": "error"})
	if !f.Match(line) {
		t.Error("expected true when value equals suffix exactly")
	}
}
