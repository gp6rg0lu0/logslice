package cli

import (
	"testing"
)

func TestParseRegionFlag_Empty(t *testing.T) {
	_, err := ParseRegionFlag("")
	if err == nil {
		t.Fatal("expected error for empty flag")
	}
}

func TestParseRegionFlag_MissingColon(t *testing.T) {
	_, err := ParseRegionFlag("regionus-east-1")
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseRegionFlag_EmptyField(t *testing.T) {
	_, err := ParseRegionFlag(":us-east-1")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestParseRegionFlag_NoRegions(t *testing.T) {
	_, err := ParseRegionFlag("region:")
	if err == nil {
		t.Fatal("expected error when no regions provided")
	}
}

func TestParseRegionFlag_Valid(t *testing.T) {
	tests := []struct {
		input    string
		field    string
		nRegions int
	}{
		{"region:us-east-1", "region", 1},
		{"region:us-east-1,eu-west-2", "region", 2},
		{"dc:us-*,eu-*", "dc", 2},
		{"region: us-east-1 , eu-west-2 ", "region", 2},
	}
	for _, tc := range tests {
		f, err := ParseRegionFlag(tc.input)
		if err != nil {
			t.Errorf("input %q: unexpected error: %v", tc.input, err)
			continue
		}
		if f.Field() != tc.field {
			t.Errorf("input %q: expected field %q, got %q", tc.input, tc.field, f.Field())
		}
		if len(f.Regions()) != tc.nRegions {
			t.Errorf("input %q: expected %d regions, got %d", tc.input, tc.nRegions, len(f.Regions()))
		}
	}
}

func TestParseRegionFlag_Errors(t *testing.T) {
	bad := []string{
		"",
		"regiononly",
		":us-east-1",
		"region:",
	}
	for _, b := range bad {
		_, err := ParseRegionFlag(b)
		if err == nil {
			t.Errorf("expected error for input %q", b)
		}
	}
}
