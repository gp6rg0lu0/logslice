package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeSessionLogLine(field, value string) *parser.LogLine {
	fields := map[string]string{}
	if field != "" && value != "" {
		fields[field] = value
	}
	return parser.NewLogLine(fields)
}

func TestNewSessionFilter_EmptyField(t *testing.T) {
	_, err := NewSessionFilter("", []string{"abc"}, false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewSessionFilter_NoValues(t *testing.T) {
	_, err := NewSessionFilter("session_id", nil, false)
	if err == nil {
		t.Fatal("expected error for no values")
	}
}

func TestNewSessionFilter_BlankValue(t *testing.T) {
	_, err := NewSessionFilter("session_id", []string{"abc", "  "}, false)
	if err == nil {
		t.Fatal("expected error for blank value")
	}
}

func TestSessionFilter_Accessors(t *testing.T) {
	f, err := NewSessionFilter("session_id", []string{"s1", "s2"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "session_id" {
		t.Errorf("Field() = %q, want %q", f.Field(), "session_id")
	}
	if len(f.Values()) != 2 {
		t.Errorf("Values() len = %d, want 2", len(f.Values()))
	}
	if !f.Prefix() {
		t.Error("Prefix() should be true")
	}
}

func TestSessionFilter_NilLine(t *testing.T) {
	f, _ := NewSessionFilter("session_id", []string{"abc"}, false)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestSessionFilter_ExactMatch(t *testing.T) {
	f, _ := NewSessionFilter("session_id", []string{"abc123", "xyz789"}, false)
	if !f.Match(makeSessionLogLine("session_id", "abc123")) {
		t.Error("expected match for exact value")
	}
	if f.Match(makeSessionLogLine("session_id", "abc")) {
		t.Error("expected no match for partial value in exact mode")
	}
}

func TestSessionFilter_PrefixMatch(t *testing.T) {
	f, _ := NewSessionFilter("session_id", []string{"sess-"}, true)
	if !f.Match(makeSessionLogLine("session_id", "sess-abc123")) {
		t.Error("expected prefix match")
	}
	if f.Match(makeSessionLogLine("session_id", "other-abc")) {
		t.Error("expected no match for non-matching prefix")
	}
}

func TestSessionFilter_MissingField(t *testing.T) {
	f, _ := NewSessionFilter("session_id", []string{"abc"}, false)
	if f.Match(makeSessionLogLine("", "")) {
		t.Error("expected no match when field is missing")
	}
}
