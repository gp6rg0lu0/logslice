package cli

import (
	"testing"
)

func TestParseSourceFlag_Empty(t *testing.T) {
	f, err := ParseSourceFlag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != nil {
		t.Error("expected nil filter for empty input")
	}
}

func TestParseSourceFlag_Valid(t *testing.T) {
	f, err := ParseSourceFlag("source:app,worker")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.Field() != "source" {
		t.Errorf("expected field 'source', got %q", f.Field())
	}
	if len(f.Sources()) != 2 {
		t.Errorf("expected 2 sources, got %d", len(f.Sources()))
	}
}

func TestParseSourceFlag_Errors(t *testing.T) {
	cases := []string{
		"noseparator",
		":nofieldname",
		"field:",
	}
	for _, c := range cases {
		_, err := ParseSourceFlag(c)
		if err == nil {
			t.Errorf("expected error for input %q", c)
		}
	}
}
