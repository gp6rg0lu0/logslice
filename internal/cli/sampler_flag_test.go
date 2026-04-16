package cli

import (
	"testing"
)

func TestParseSamplerFlag_Valid(t *testing.T) {
	cases := []struct {
		input string
		want  int
	}{
		{"1", 1},
		{"5", 5},
		{"every:3", 3},
		{"every:10", 10},
		{" 4 ", 4},
	}
	for _, tc := range cases {
		got, err := ParseSamplerFlag(tc.input)
		if err != nil {
			t.Errorf("input %q: unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("input %q: expected %d, got %d", tc.input, tc.want, got)
		}
	}
}

func TestParseSamplerFlag_Errors(t *testing.T) {
	cases := []string{
		"",
		"abc",
		"0",
		"-1",
		"every:0",
		"every:abc",
	}
	for _, input := range cases {
		_, err := ParseSamplerFlag(input)
		if err == nil {
			t.Errorf("input %q: expected error, got nil", input)
		}
	}
}
