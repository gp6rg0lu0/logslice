package cli

import "testing"

func TestParseRangeFlag_Valid(t *testing.T) {
	tests := []struct {
		input string
		field string
		min   float64
		max   float64
	}{
		{"latency:10:200", "latency", 10, 200},
		{"status:200:599", "status", 200, 599},
		{"score:0.5:1.0", "score", 0.5, 1.0},
		{"count:0:0", "count", 0, 0},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			rv, err := ParseRangeFlag(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if rv.Field != tc.field {
				t.Errorf("Field = %q, want %q", rv.Field, tc.field)
			}
			if rv.Min != tc.min {
				t.Errorf("Min = %v, want %v", rv.Min, tc.min)
			}
			if rv.Max != tc.max {
				t.Errorf("Max = %v, want %v", rv.Max, tc.max)
			}
		})
	}
}

func TestParseRangeFlag_Errors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"missing parts", "latency:10"},
		{"empty field", ":10:200"},
		{"bad min", "latency:abc:200"},
		{"bad max", "latency:10:xyz"},
		{"min gt max", "latency:500:100"},
		{"no separators", "latency"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseRangeFlag(tc.input)
			if err == nil {
				t.Errorf("expected error for input %q", tc.input)
			}
		})
	}
}
