package cli

import (
	"testing"
)

func TestParseJSONPathFlag_Empty(t *testing.T) {
	filters, err := ParseJSONPathFlag(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Errorf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseJSONPathFlag_Valid(t *testing.T) {
	filters, err := ParseJSONPathFlag([]string{"meta.region=us-east-1", "env=production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
	if filters[0].Path() != "meta.region" {
		t.Errorf("expected path meta.region, got %s", filters[0].Path())
	}
	if filters[0].Value() != "us-east-1" {
		t.Errorf("expected value us-east-1, got %s", filters[0].Value())
	}
	if filters[1].Path() != "env" {
		t.Errorf("expected path env, got %s", filters[1].Path())
	}
}

func TestParseJSONPathFlag_Errors(t *testing.T) {
	cases := []struct {
		arg string
	}{
		{"noequals"},
		{"=nopath"},
		{"path="},
		{"a..b=val"},
	}
	for _, tc := range cases {
		_, err := ParseJSONPathFlag([]string{tc.arg})
		if err == nil {
			t.Errorf("expected error for arg %q", tc.arg)
		}
	}
}
