package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeRequestIDLogLine(field, value string) *parser.LogLine {
	fields := map[string]string{}
	if field != "" {
		fields[field] = value
	}
	return parser.NewLogLine(fields)
}

func TestNewRequestIDFilter_EmptyField(t *testing.T) {
	_, err := NewRequestIDFilter("", []string{"abc"}, true)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewRequestIDFilter_NoValues(t *testing.T) {
	_, err := NewRequestIDFilter("request_id", nil, true)
	if err == nil {
		t.Fatal("expected error for nil values")
	}
}

func TestNewRequestIDFilter_BlankValue(t *testing.T) {
	_, err := NewRequestIDFilter("request_id", []string{"abc", ""}, true)
	if err == nil {
		t.Fatal("expected error for blank value")
	}
}

func TestRequestIDFilter_Accessors(t *testing.T) {
	f, err := NewRequestIDFilter("request_id", []string{"req-123"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "request_id" {
		t.Errorf("Field() = %q, want %q", f.Field(), "request_id")
	}
	if len(f.Values()) != 1 || f.Values()[0] != "req-123" {
		t.Errorf("Values() = %v, want [req-123]", f.Values())
	}
	if f.Exact() {
		t.Error("Exact() should be false")
	}
}

func TestRequestIDFilter_NilLine(t *testing.T) {
	f, _ := NewRequestIDFilter("request_id", []string{"req-1"}, true)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestRequestIDFilter_ExactMatch(t *testing.T) {
	f, _ := NewRequestIDFilter("request_id", []string{"req-abc-123"}, true)
	if !f.Match(makeRequestIDLogLine("request_id", "req-abc-123")) {
		t.Error("expected match for exact value")
	}
	if f.Match(makeRequestIDLogLine("request_id", "req-abc-123-extra")) {
		t.Error("expected no match for longer value in exact mode")
	}
}

func TestRequestIDFilter_PrefixMatch(t *testing.T) {
	f, _ := NewRequestIDFilter("request_id", []string{"req-abc"}, false)
	if !f.Match(makeRequestIDLogLine("request_id", "req-abc-123")) {
		t.Error("expected prefix match")
	}
	if f.Match(makeRequestIDLogLine("request_id", "other-req-abc")) {
		t.Error("expected no match when prefix not at start")
	}
}

func TestRequestIDFilter_MissingField(t *testing.T) {
	f, _ := NewRequestIDFilter("request_id", []string{"req-1"}, true)
	if f.Match(makeRequestIDLogLine("", "req-1")) {
		t.Error("expected false when field is missing")
	}
}

func TestRequestIDFilter_MultipleValues(t *testing.T) {
	f, _ := NewRequestIDFilter("request_id", []string{"req-1", "req-2", "req-3"}, true)
	for _, id := range []string{"req-1", "req-2", "req-3"} {
		if !f.Match(makeRequestIDLogLine("request_id", id)) {
			t.Errorf("expected match for %q", id)
		}
	}
	if f.Match(makeRequestIDLogLine("request_id", "req-4")) {
		t.Error("expected no match for req-4")
	}
}
