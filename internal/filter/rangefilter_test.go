package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func makeRangeLogLine(fields map[string]string) *parser.LogLine {
	raw := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		raw[k] = v
	}
	return parser.NewLogLine(raw)
}

func TestNewRangeFilter_EmptyField(t *testing.T) {
	_, err := filter.NewRangeFilter("", 0, 10)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewRangeFilter_InvalidRange(t *testing.T) {
	_, err := filter.NewRangeFilter("latency", 100, 10)
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestRangeFilter_Accessors(t *testing.T) {
	f, err := filter.NewRangeFilter("latency", 5, 50)
	if err != nil {
		t.Fatal(err)
	}
	if f.Field() != "latency" {
		t.Errorf("Field() = %q, want %q", f.Field(), "latency")
	}
	if f.Min() != 5 {
		t.Errorf("Min() = %v, want 5", f.Min())
	}
	if f.Max() != 50 {
		t.Errorf("Max() = %v, want 50", f.Max())
	}
}

func TestRangeFilter_Match(t *testing.T) {
	f, _ := filter.NewRangeFilter("latency", 10, 100)
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"below range", "5", false},
		{"at min", "10", true},
		{"in range", "55", true},
		{"at max", "100", true},
		{"above range", "200", false},
		{"non-numeric", "fast", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			line := makeRangeLogLine(map[string]string{"latency": tc.value})
			if got := f.Match(line); got != tc.want {
				t.Errorf("Match() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRangeFilter_MissingField(t *testing.T) {
	f, _ := filter.NewRangeFilter("latency", 0, 100)
	line := makeRangeLogLine(map[string]string{"other": "42"})
	if f.Match(line) {
		t.Error("expected false for missing field")
	}
}

func TestRangeFilter_NilLine(t *testing.T) {
	f, _ := filter.NewRangeFilter("latency", 0, 100)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}
