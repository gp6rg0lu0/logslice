package cli

import (
	"testing"
)

func TestParseRequestIDFlag_Empty(t *testing.T) {
	f, err := ParseRequestIDFlag("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != nil {
		t.Error("expected nil filter for empty input")
	}
}

func TestParseRequestIDFlag_Valid(t *testing.T) {
	tests := []struct {
		input   string
		exact   bool
		field   string
		nValues int
	}{
		{"request_id:exact:req-abc", true, "request_id", 1},
		{"req_id:prefix:req-1,req-2,req-3", false, "req_id", 3},
		{"x_request_id:exact:abc-123,def-456", true, "x_request_id", 2},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			f, err := ParseRequestIDFlag(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f == nil {
				t.Fatal("expected non-nil filter")
			}
			if f.Field() != tc.field {
				t.Errorf("Field() = %q, want %q", f.Field(), tc.field)
			}
			if f.Exact() != tc.exact {
				t.Errorf("Exact() = %v, want %v", f.Exact(), tc.exact)
			}
			if len(f.Values()) != tc.nValues {
				t.Errorf("len(Values()) = %d, want %d", len(f.Values()), tc.nValues)
			}
		})
	}
}

func TestParseRequestIDFlag_Errors(t *testing.T) {
	cases := []string{
		"request_id:exact",          // missing IDs segment
		"request_id",                // missing mode and IDs
		":exact:req-1",              // empty field
		"request_id:badmode:req-1",  // invalid mode
		"request_id:exact:",         // empty IDs
	}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			_, err := ParseRequestIDFlag(c)
			if err == nil {
				t.Errorf("expected error for input %q", c)
			}
		})
	}
}
