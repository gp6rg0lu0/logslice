package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeLenLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]string{field: value})
}

func TestNewLenFilter_EmptyField(t *testing.T) {
	_, err := NewLenFilter("", 0, -1)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewLenFilter_NegativeMin(t *testing.T) {
	_, err := NewLenFilter("msg", -1, -1)
	if err == nil {
		t.Fatal("expected error for negative minLen")
	}
}

func TestNewLenFilter_MaxLessThanMin(t *testing.T) {
	_, err := NewLenFilter("msg", 10, 5)
	if err == nil {
		t.Fatal("expected error when maxLen < minLen")
	}
}

func TestLenFilter_Accessors(t *testing.T) {
	f, err := NewLenFilter("msg", 2, 8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "msg" {
		t.Errorf("expected field 'msg', got %q", f.Field())
	}
	if f.MinLen() != 2 {
		t.Errorf("expected minLen 2, got %d", f.MinLen())
	}
	if f.MaxLen() != 8 {
		t.Errorf("expected maxLen 8, got %d", f.MaxLen())
	}
}

func TestLenFilter_NilLine(t *testing.T) {
	f, _ := NewLenFilter("msg", 0, -1)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestLenFilter_MissingField(t *testing.T) {
	f, _ := NewLenFilter("msg", 0, -1)
	line := makeLenLogLine("other", "hello")
	if f.Match(line) {
		t.Error("expected false for missing field")
	}
}

func TestLenFilter_Match(t *testing.T) {
	tests := []struct {
		value  string
		min    int
		max    int
		match  bool
	}{
		{"hello", 3, 10, true},
		{"hi", 3, 10, false},
		{"hello world!", 3, 10, false},
		{"hello", 5, 5, true},
		{"hello", 0, -1, true},
		{"", 0, 0, true},
		{"", 1, -1, false},
	}
	for _, tt := range tests {
		f, err := NewLenFilter("msg", tt.min, tt.max)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		line := makeLenLogLine("msg", tt.value)
		if got := f.Match(line); got != tt.match {
			t.Errorf("value=%q min=%d max=%d: expected %v, got %v", tt.value, tt.min, tt.max, tt.match, got)
		}
	}
}
