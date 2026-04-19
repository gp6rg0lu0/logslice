package filter_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeDurationLogLine(field, value string) *parser.LogLine {
	fields := map[string]any{}
	if field != "" {
		fields[field] = value
	}
	return parser.NewLogLine(fields)
}

func TestNewDurationFilter_EmptyField(t *testing.T) {
	_, err := filter.NewDurationFilter("", "1s", "5s")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewDurationFilter_InvalidMin(t *testing.T) {
	_, err := filter.NewDurationFilter("latency", "bad", "5s")
	if err == nil {
		t.Fatal("expected error for invalid min")
	}
}

func TestNewDurationFilter_InvalidMax(t *testing.T) {
	_, err := filter.NewDurationFilter("latency", "1s", "bad")
	if err == nil {
		t.Fatal("expected error for invalid max")
	}
}

func TestNewDurationFilter_MaxLessThanMin(t *testing.T) {
	_, err := filter.NewDurationFilter("latency", "10s", "1s")
	if err == nil {
		t.Fatal("expected error when max < min")
	}
}

func TestDurationFilter_Accessors(t *testing.T) {
	f, err := filter.NewDurationFilter("latency", "1s", "10s")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "latency" {
		t.Errorf("expected field latency, got %s", f.Field())
	}
	if f.Min() != time.Second {
		t.Errorf("expected min 1s, got %v", f.Min())
	}
	if f.Max() != 10*time.Second {
		t.Errorf("expected max 10s, got %v", f.Max())
	}
}

func TestDurationFilter_NilLine(t *testing.T) {
	f, _ := filter.NewDurationFilter("latency", "1s", "5s")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestDurationFilter_Match(t *testing.T) {
	f, _ := filter.NewDurationFilter("latency", "1s", "5s")
	cases := []struct {
		val  string
		want bool
	}{
		{"500ms", false},
		{"1s", true},
		{"3s", true},
		{"5s", true},
		{"6s", false},
		{"notaduration", false},
	}
	for _, tc := range cases {
		line := makeDurationLogLine("latency", tc.val)
		if got := f.Match(line); got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.val, got, tc.want)
		}
	}
}

func TestDurationFilter_MissingField(t *testing.T) {
	f, _ := filter.NewDurationFilter("latency", "1s", "5s")
	line := makeDurationLogLine("", "")
	if f.Match(line) {
		t.Error("expected false for missing field")
	}
}
