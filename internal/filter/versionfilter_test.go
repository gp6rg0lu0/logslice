package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeVersionLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]string{field: value})
}

func TestNewVersionFilter_EmptyField(t *testing.T) {
	_, err := NewVersionFilter("", "1.0.0", "")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewVersionFilter_NoBounds(t *testing.T) {
	_, err := NewVersionFilter("version", "", "")
	if err == nil {
		t.Fatal("expected error when no bounds set")
	}
}

func TestNewVersionFilter_InvalidMin(t *testing.T) {
	_, err := NewVersionFilter("version", "abc", "")
	if err == nil {
		t.Fatal("expected error for invalid min")
	}
}

func TestNewVersionFilter_InvalidMax(t *testing.T) {
	_, err := NewVersionFilter("version", "", "1.x.0")
	if err == nil {
		t.Fatal("expected error for invalid max")
	}
}

func TestNewVersionFilter_MinExceedsMax(t *testing.T) {
	_, err := NewVersionFilter("version", "2.0.0", "1.0.0")
	if err == nil {
		t.Fatal("expected error when min exceeds max")
	}
}

func TestVersionFilter_Accessors(t *testing.T) {
	f, err := NewVersionFilter("version", "1.0.0", "2.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "version" {
		t.Errorf("expected field 'version', got %q", f.Field())
	}
}

func TestVersionFilter_NilLine(t *testing.T) {
	f, _ := NewVersionFilter("version", "1.0.0", "")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestVersionFilter_Match(t *testing.T) {
	tests := []struct {
		name    string
		min     string
		max     string
		value   string
		want    bool
	}{
		{"within range", "1.0.0", "2.0.0", "1.5.0", true},
		{"at min", "1.0.0", "2.0.0", "1.0.0", true},
		{"at max", "1.0.0", "2.0.0", "2.0.0", true},
		{"below min", "1.0.0", "2.0.0", "0.9.9", false},
		{"above max", "1.0.0", "2.0.0", "2.0.1", false},
		{"only min set", "1.2.0", "", "1.3.0", true},
		{"only min fails", "1.2.0", "", "1.1.0", false},
		{"only max set", "", "3.0.0", "2.9.9", true},
		{"only max fails", "", "3.0.0", "3.0.1", false},
		{"invalid value", "1.0.0", "2.0.0", "not-a-version", false},
		{"missing field", "1.0.0", "2.0.0", "", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := NewVersionFilter("ver", tc.min, tc.max)
			if err != nil {
				t.Fatalf("setup error: %v", err)
			}
			var line *parser.LogLine
			if tc.value != "" {
				line = makeVersionLogLine("ver", tc.value)
			} else {
				line = parser.NewLogLine(map[string]string{})
			}
			if got := f.Match(line); got != tc.want {
				t.Errorf("Match() = %v, want %v", got, tc.want)
			}
		})
	}
}
