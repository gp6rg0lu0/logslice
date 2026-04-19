package cli

import (
	"testing"
)

func TestParseIPFlag_Empty(t *testing.T) {
	filters, err := ParseIPFlag(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 0 {
		t.Errorf("expected 0 filters, got %d", len(filters))
	}
}

func TestParseIPFlag_Valid(t *testing.T) {
	filters, err := ParseIPFlag([]string{"src=10.0.0.0/8", "dst=192.168.0.0/16"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}
	if filters[0].Field() != "src" {
		t.Errorf("expected field src, got %s", filters[0].Field())
	}
	if filters[1].CIDR() != "192.168.0.0/16" {
		t.Errorf("expected cidr 192.168.0.0/16, got %s", filters[1].CIDR())
	}
}

func TestParseIPFlag_Errors(t *testing.T) {
	cases := []string{
		"noequals",
		"=10.0.0.0/8",
		"field=",
		"field=not-a-cidr",
	}
	for _, c := range cases {
		_, err := ParseIPFlag([]string{c})
		if err == nil {
			t.Errorf("expected error for input %q", c)
		}
	}
}
