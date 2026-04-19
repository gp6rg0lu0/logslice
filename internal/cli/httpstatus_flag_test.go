package cli_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/cli"
)

func TestParseHTTPStatusFlag_Empty(t *testing.T) {
	f, err := cli.ParseHTTPStatusFlag("")
	if err != nil {
		t.Fatal(err)
	}
	if f != nil {
		t.Fatal("expected nil filter for empty flag")
	}
}

func TestParseHTTPStatusFlag_Valid(t *testing.T) {
	cases := []struct {
		input    string
		field    string
		min, max int
	}{
		{"status:200-299", "status", 200, 299},
		{"http_code:400-499", "http_code", 400, 499},
		{"code:500-599", "code", 500, 599},
	}
	for _, c := range cases {
		f, err := cli.ParseHTTPStatusFlag(c.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", c.input, err)
		}
		if f.Field() != c.field {
			t.Errorf("field: got %q, want %q", f.Field(), c.field)
		}
		if f.Min() != c.min || f.Max() != c.max {
			t.Errorf("bounds: got %d-%d, want %d-%d", f.Min(), f.Max(), c.min, c.max)
		}
	}
}

func TestParseHTTPStatusFlag_Errors(t *testing.T) {
	cases := []string{
		"status",       // no colon
		":200-299",     // empty field
		"status:200",   // no dash
		"status:abc-299",
		"status:200-xyz",
		"status:50-200",  // min out of range
		"status:400-300", // min > max
	}
	for _, c := range cases {
		_, err := cli.ParseHTTPStatusFlag(c)
		if err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}
