package cli_test

import (
	"testing"

	"github.com/nicholasgasior/logslice/internal/cli"
)

func TestParseTypeFlag_Empty(t *testing.T) {
	f, err := cli.ParseTypeFlag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != nil {
		t.Error("expected nil filter for empty input")
	}
}

func TestParseTypeFlag_Valid(t *testing.T) {
	tests := []struct {
		input    string
		field    string
		wantType string
	}{
		{"msg:string", "msg", "string"},
		{"count:number", "count", "number"},
		{"ok:bool", "ok", "bool"},
		{"meta:object", "meta", "object"},
		{"tags:array", "tags", "array"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			f, err := cli.ParseTypeFlag(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f == nil {
				t.Fatal("expected non-nil filter")
			}
			if f.Field() != tt.field || f.WantType() != tt.wantType {
				t.Errorf("got field=%q type=%q, want field=%q type=%q", f.Field(), f.WantType(), tt.field, tt.wantType)
			}
		})
	}
}

func TestParseTypeFlag_Errors(t *testing.T) {
	bad := []string{
		"nocodon",
		":string",
		"field:",
		"field:integer",
	}
	for _, s := range bad {
		t.Run(s, func(t *testing.T) {
			_, err := cli.ParseTypeFlag(s)
			if err == nil {
				t.Errorf("expected error for input %q", s)
			}
		})
	}
}
