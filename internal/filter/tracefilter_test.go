package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeTraceLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]interface{}{field: value})
}

func TestNewTraceFilter_EmptyField(t *testing.T) {
	_, err := NewTraceFilter("", []string{"abc"}, true)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewTraceFilter_NoValues(t *testing.T) {
	_, err := NewTraceFilter("trace_id", nil, true)
	if err == nil {
		t.Fatal("expected error for empty values")
	}
}

func TestNewTraceFilter_BlankValue(t *testing.T) {
	_, err := NewTraceFilter("trace_id", []string{"  "}, true)
	if err == nil {
		t.Fatal("expected error for blank value")
	}
}

func TestTraceFilter_Accessors(t *testing.T) {
	f, err := NewTraceFilter("trace_id", []string{"abc123"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "trace_id" {
		t.Errorf("expected trace_id, got %s", f.Field())
	}
	if f.Exact() {
		t.Error("expected exact=false")
	}
	if len(f.Values()) != 1 || f.Values()[0] != "abc123" {
		t.Errorf("unexpected values: %v", f.Values())
	}
}

func TestTraceFilter_NilLine(t *testing.T) {
	f, _ := NewTraceFilter("trace_id", []string{"abc"}, true)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestTraceFilter_ExactMatch(t *testing.T) {
	f, _ := NewTraceFilter("trace_id", []string{"abc123", "def456"}, true)
	if !f.Match(makeTraceLogLine("trace_id", "abc123")) {
		t.Error("expected match for abc123")
	}
	if f.Match(makeTraceLogLine("trace_id", "abc")) {
		t.Error("expected no match for prefix-only value in exact mode")
	}
}

func TestTraceFilter_PrefixMatch(t *testing.T) {
	f, _ := NewTraceFilter("trace_id", []string{"abc"}, false)
	if !f.Match(makeTraceLogLine("trace_id", "abc123")) {
		t.Error("expected prefix match")
	}
	if f.Match(makeTraceLogLine("trace_id", "xyz999")) {
		t.Error("expected no match for non-prefix value")
	}
}

func TestTraceFilter_MissingField(t *testing.T) {
	f, _ := NewTraceFilter("trace_id", []string{"abc"}, true)
	line := parser.NewLogLine(map[string]interface{}{"other": "abc"})
	if f.Match(line) {
		t.Error("expected no match when field is missing")
	}
}
