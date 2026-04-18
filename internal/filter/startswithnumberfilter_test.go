package filter_test

import (
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeStartsWithNumberLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]string{field: value})
}

func TestNewStartsWithNumberFilter_EmptyField(t *testing.T) {
	_, err := filter.NewStartsWithNumberFilter("", false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewStartsWithNumberFilter_Valid(t *testing.T) {
	f, err := filter.NewStartsWithNumberFilter("code", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "code" {
		t.Errorf("expected field 'code', got %q", f.Field())
	}
	if f.Invert() {
		t.Error("expected invert=false")
	}
}

func TestStartsWithNumberFilter_NilLine(t *testing.T) {
	f, _ := filter.NewStartsWithNumberFilter("code", false)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestStartsWithNumberFilter_MatchesDigitStart(t *testing.T) {
	f, _ := filter.NewStartsWithNumberFilter("msg", false)
	line := makeStartsWithNumberLogLine("msg", "404 not found")
	if !f.Match(line) {
		t.Error("expected match for value starting with digit")
	}
}

func TestStartsWithNumberFilter_NoMatchLetterStart(t *testing.T) {
	f, _ := filter.NewStartsWithNumberFilter("msg", false)
	line := makeStartsWithNumberLogLine("msg", "error occurred")
	if f.Match(line) {
		t.Error("expected no match for value starting with letter")
	}
}

func TestStartsWithNumberFilter_InvertedMatch(t *testing.T) {
	f, _ := filter.NewStartsWithNumberFilter("msg", true)
	line := makeStartsWithNumberLogLine("msg", "error occurred")
	if !f.Match(line) {
		t.Error("expected match when inverted and value starts with letter")
	}
}

func TestStartsWithNumberFilter_MissingField(t *testing.T) {
	f, _ := filter.NewStartsWithNumberFilter("code", false)
	line := makeStartsWithNumberLogLine("other", "123")
	if f.Match(line) {
		t.Error("expected no match for missing field")
	}
}

func TestStartsWithNumberFilter_MissingFieldInverted(t *testing.T) {
	f, _ := filter.NewStartsWithNumberFilter("code", true)
	line := makeStartsWithNumberLogLine("other", "123")
	if !f.Match(line) {
		t.Error("expected match when inverted and field is missing")
	}
}
