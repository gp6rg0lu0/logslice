package filter_test

import (
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeBoolLogLine(fields map[string]any) *parser.LogLine {
	return parser.NewLogLine(fields)
}

func TestNewBoolFilter_EmptyField(t *testing.T) {
	_, err := filter.NewBoolFilter("", "true")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewBoolFilter_InvalidValue(t *testing.T) {
	_, err := filter.NewBoolFilter("active", "yes")
	if err == nil {
		t.Fatal("expected error for invalid value")
	}
}

func TestNewBoolFilter_Valid(t *testing.T) {
	f, err := filter.NewBoolFilter("active", "true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "active" {
		t.Errorf("expected field 'active', got %q", f.Field())
	}
	if !f.Want() {
		t.Error("expected Want() to be true")
	}
}

func TestBoolFilter_NilLine(t *testing.T) {
	f, _ := filter.NewBoolFilter("active", "true")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestBoolFilter_MissingField(t *testing.T) {
	f, _ := filter.NewBoolFilter("active", "true")
	line := makeBoolLogLine(map[string]any{"msg": "hello"})
	if f.Match(line) {
		t.Error("expected false for missing field")
	}
}

func TestBoolFilter_MatchesBoolTrue(t *testing.T) {
	f, _ := filter.NewBoolFilter("active", "true")
	line := makeBoolLogLine(map[string]any{"active": true})
	if !f.Match(line) {
		t.Error("expected match for bool true")
	}
}

func TestBoolFilter_NoMatchBoolFalse(t *testing.T) {
	f, _ := filter.NewBoolFilter("active", "true")
	line := makeBoolLogLine(map[string]any{"active": false})
	if f.Match(line) {
		t.Error("expected no match for bool false")
	}
}

func TestBoolFilter_MatchesStringTrue(t *testing.T) {
	f, _ := filter.NewBoolFilter("active", "true")
	line := makeBoolLogLine(map[string]any{"active": "true"})
	if !f.Match(line) {
		t.Error("expected match for string 'true'")
	}
}

func TestBoolFilter_FalseValue(t *testing.T) {
	f, _ := filter.NewBoolFilter("active", "false")
	line := makeBoolLogLine(map[string]any{"active": false})
	if !f.Match(line) {
		t.Error("expected match for bool false with want=false")
	}
}
