package cli_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/cli"
)

func TestParseTagFlag_Empty(t *testing.T) {
	_, err := cli.ParseTagFlag("")
	if err == nil {
		t.Fatal("expected error for empty input")
	}
}

func TestParseTagFlag_MissingColon(t *testing.T) {
	_, err := cli.ParseTagFlag("justfield")
	if err == nil {
		t.Fatal("expected error when colon separator is absent")
	}
}

func TestParseTagFlag_EmptyField(t *testing.T) {
	_, err := cli.ParseTagFlag(":prod,staging")
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestParseTagFlag_Valid(t *testing.T) {
	tests := []struct {
		input     string
		field     string
		numTags   int
	}{
		{"tags:prod,staging", "tags", 2},
		{"labels:canary", "labels", 1},
		{"env:prod,qa,dev", "env", 3},
	}
	for _, tc := range tests {
		f, err := cli.ParseTagFlag(tc.input)
		if err != nil {
			t.Errorf("input %q: unexpected error: %v", tc.input, err)
			continue
		}
		if f.Field() != tc.field {
			t.Errorf("input %q: expected field %q, got %q", tc.input, tc.field, f.Field())
		}
		if len(f.Tags()) != tc.numTags {
			t.Errorf("input %q: expected %d tags, got %d", tc.input, tc.numTags, len(f.Tags()))
		}
	}
}

func TestParseTagFlag_Errors(t *testing.T) {
	inputs := []string{
		"tags:",      // no tags after colon
		"tags:  , ", // only blank tags
	}
	for _, raw := range inputs {
		_, err := cli.ParseTagFlag(raw)
		if err == nil {
			t.Errorf("expected error for input %q", raw)
		}
	}
}
