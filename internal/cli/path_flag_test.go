package cli

import (
	"testing"
)

func TestParsePathFlag_Empty(t *testing.T) {
	_, err := ParsePathFlag("")
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestParsePathFlag_MissingColon(t *testing.T) {
	_, err := ParsePathFlag("pathfield")
	if err == nil {
		t.Fatal("expected error when colon is missing")
	}
}

func TestParsePathFlag_Valid_PrefixMode(t *testing.T) {
	f, err := ParsePathFlag("url:/api/v1,/health")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "url" {
		t.Errorf("expected field 'url', got %q", f.Field())
	}
	if f.Exact() {
		t.Error("expected prefix mode (exact=false)")
	}
	if len(f.Paths()) != 2 {
		t.Errorf("expected 2 paths, got %d", len(f.Paths()))
	}
}

func TestParsePathFlag_Valid_ExactMode(t *testing.T) {
	f, err := ParsePathFlag("path:exact:/ready,/healthz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "path" {
		t.Errorf("expected field 'path', got %q", f.Field())
	}
	if !f.Exact() {
		t.Error("expected exact=true")
	}
	if len(f.Paths()) != 2 {
		t.Errorf("expected 2 paths, got %d", len(f.Paths()))
	}
}

func TestParsePathFlag_Errors(t *testing.T) {
	cases := []struct {
		input string
		desc  string
	}{
		{":  , ", "blank paths only"},
		{"field:", "empty path list"},
	}
	for _, tc := range cases {
		_, err := ParsePathFlag(tc.input)
		if err == nil {
			t.Errorf("expected error for %s (input=%q)", tc.desc, tc.input)
		}
	}
}
