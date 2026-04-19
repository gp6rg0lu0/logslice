package cli

import (
	"testing"
)

func TestParseTraceFlag_Empty(t *testing.T) {
	_, err := ParseTraceFlag("")
	if err == nil {
		t.Fatal("expected error for empty flag")
	}
}

func TestParseTraceFlag_MissingColon(t *testing.T) {
	_, err := ParseTraceFlag("trace_id")
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseTraceFlag_Valid_PrefixMode(t *testing.T) {
	f, err := ParseTraceFlag("trace_id:abc,def")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "trace_id" {
		t.Errorf("expected trace_id, got %s", f.Field())
	}
	if f.Exact() {
		t.Error("expected prefix mode")
	}
	if len(f.Values()) != 2 {
		t.Errorf("expected 2 values, got %d", len(f.Values()))
	}
}

func TestParseTraceFlag_Valid_ExactMode(t *testing.T) {
	f, err := ParseTraceFlag("trace_id:abc123,exact")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Exact() {
		t.Error("expected exact mode")
	}
	if len(f.Values()) != 1 || f.Values()[0] != "abc123" {
		t.Errorf("unexpected values: %v", f.Values())
	}
}

func TestParseTraceFlag_Errors(t *testing.T) {
	cases := []string{
		"trace_id:",
		"trace_id:  ,  ",
	}
	for _, c := range cases {
		_, err := ParseTraceFlag(c)
		if err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}
