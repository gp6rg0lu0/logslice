package cli_test

import (
	"testing"

	"github.com/user/logslice/internal/cli"
)

func TestParseContainsFlag_Empty(t *testing.T) {
	f, err := cli.ParseContainsFlag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != nil {
		t.Error("expected nil filter for empty input")
	}
}

func TestParseContainsFlag_Valid(t *testing.T) {
	tests := []struct {
		input    string
		field    string
		sub      string
		fold     bool
	}{
		{"msg:hello", "msg", "hello", false},
		{"msg:hello:false", "msg", "hello", false},
		{"msg:hello:true", "msg", "hello", true},
		{"msg:hello:fold", "msg", "hello", true},
		{"level:err:1", "level", "err", true},
	}
	for _, tt := range tests {
		f, err := cli.ParseContainsFlag(tt.input)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", tt.input, err)
			continue
		}
		if f.Field() != tt.field {
			t.Errorf("%q: Field() = %q, want %q", tt.input, f.Field(), tt.field)
		}
		if f.Substring() != tt.sub {
			t.Errorf("%q: Substring() = %q, want %q", tt.input, f.Substring(), tt.sub)
		}
		if f.CaseFold() != tt.fold {
			t.Errorf("%q: CaseFold() = %v, want %v", tt.input, f.CaseFold(), tt.fold)
		}
	}
}

func TestParseContainsFlag_Errors(t *testing.T) {
	tests := []string{
		"onlyfield",
		"msg:hello:badval",
		"msg::true",
		":sub",
	}
	for _, input := range tests {
		_, err := cli.ParseContainsFlag(input)
		if err == nil {
			t.Errorf("%q: expected error, got nil", input)
		}
	}
}
