package filter_test

import (
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeBetweenLogLine(fields map[string]string) *parser.LogLine {
	raw := make(map[string]interface{})
	for k, v := range fields {
		raw[k] = v
	}
	return parser.NewLogLine(raw)
}

func TestNewBetweenFilter_EmptyField(t *testing.T) {
	_, err := filter.NewBetweenFilter("", 0, 10)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewBetweenFilter_MinExceedsMax(t *testing.T) {
	_, err := filter.NewBetweenFilter("latency", 100, 10)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestBetweenFilter_Accessors(t *testing.T) {
	f, err := filter.NewBetweenFilter("latency", 5, 50)
	if err != nil {
		t.Fatal(err)
	}
	if f.Field() != "latency" {
		t.Errorf("expected latency, got %s", f.Field())
	}
	if f.Min() != 5 {
		t.Errorf("expected min 5, got %v", f.Min())
	}
	if f.Max() != 50 {
		t.Errorf("expected max 50, got %v", f.Max())
	}
}

func TestBetweenFilter_NilLine(t *testing.T) {
	f, _ := filter.NewBetweenFilter("latency", 0, 100)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestBetweenFilter_Match(t *testing.T) {
	f, _ := filter.NewBetweenFilter("latency", 10, 100)
	tests := []struct {
		val  string
		want bool
	}{
		{"10", true},
		{"55", true},
		{"100", true},
		{"9", false},
		{"101", false},
		{"abc", false},
		{"", false},
	}
	for _, tt := range tests {
		line := makeBetweenLogLine(map[string]string{"latency": tt.val})
		if got := f.Match(line); got != tt.want {
			t.Errorf("val=%q: expected %v, got %v", tt.val, tt.want, got)
		}
	}
}

func TestBetweenFilter_MissingField(t *testing.T) {
	f, _ := filter.NewBetweenFilter("latency", 0, 100)
	line := makeBetweenLogLine(map[string]string{"other": "50"})
	if f.Match(line) {
		t.Error("expected false for missing field")
	}
}
