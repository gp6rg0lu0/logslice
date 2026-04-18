package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEndsWithNumberLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]interface{}{field: value})
}

func TestNewEndsWithNumberFilter_EmptyField(t *testing.T) {
	_, err := NewEndsWithNumberFilter("", false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewEndsWithNumberFilter_Valid(t *testing.T) {
	f, err := NewEndsWithNumberFilter("msg", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "msg" {
		t.Errorf("expected field 'msg', got %q", f.Field())
	}
	if f.Invert() {
		t.Error("expected invert=false")
	}
}

func TestEndsWithNumberFilter_NilLine(t *testing.T) {
	f, _ := NewEndsWithNumberFilter("msg", false)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestEndsWithNumberFilter_MatchesDigitEnd(t *testing.T) {
	f, _ := NewEndsWithNumberFilter("msg", false)
	tests := []struct {
		val  string
		want bool
	}{
		{"error42", true},
		{"event0", true},
		{"hello", false},
		{"12abc", false},
		{"9", true},
	}
	for _, tt := range tests {
		line := makeEndsWithNumberLogLine("msg", tt.val)
		if got := f.Match(line); got != tt.want {
			t.Errorf("Match(%q) = %v, want %v", tt.val, got, tt.want)
		}
	}
}

func TestEndsWithNumberFilter_Inverted(t *testing.T) {
	f, _ := NewEndsWithNumberFilter("msg", true)
	line := makeEndsWithNumberLogLine("msg", "event3")
	if f.Match(line) {
		t.Error("expected false for inverted match on digit-ending value")
	}
	line2 := makeEndsWithNumberLogLine("msg", "hello")
	if !f.Match(line2) {
		t.Error("expected true for inverted match on non-digit-ending value")
	}
}

func TestEndsWithNumberFilter_MissingField(t *testing.T) {
	f, _ := NewEndsWithNumberFilter("msg", false)
	line := makeEndsWithNumberLogLine("other", "value9")
	if f.Match(line) {
		t.Error("expected false when field is missing")
	}
}
