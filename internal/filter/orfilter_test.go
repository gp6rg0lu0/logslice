package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func makeOrLogLine(fields map[string]string) *parser.LogLine {
	raw := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		raw[k] = v
	}
	return parser.NewLogLine(raw)
}

func TestOrFilter_NoFilters_NeverMatches(t *testing.T) {
	or_ := filter.NewOrFilter()
	line := makeOrLogLine(map[string]string{"level": "info"})
	if or_.Match(line) {
		t.Error("empty OrFilter should never match")
	}
}

func TestOrFilter_IgnoresNilFilters(t *testing.T) {
	or_ := filter.NewOrFilter(nil, nil)
	if or_.Len() != 0 {
		t.Errorf("expected 0 active filters, got %d", or_.Len())
	}
}

func TestOrFilter_MatchesWhenOneMatches(t *testing.T) {
	f1, _ := filter.NewLevelFilter("debug")
	f2, _ := filter.NewLevelFilter("error")
	or_ := filter.NewOrFilter(f1, f2)

	tests := []struct {
		level string
		want  bool
	}{
		{"debug", true},
		{"error", true},
		{"info", false},
		{"warn", false},
	}
	for _, tt := range tests {
		line := makeOrLogLine(map[string]string{"level": tt.level})
		if got := or_.Match(line); got != tt.want {
			t.Errorf("level=%q: got %v, want %v", tt.level, got, tt.want)
		}
	}
}

func TestOrFilter_Len(t *testing.T) {
	f1, _ := filter.NewLevelFilter("info")
	f2, _ := filter.NewLevelFilter("warn")
	or_ := filter.NewOrFilter(f1, nil, f2)
	if or_.Len() != 2 {
		t.Errorf("expected Len 2, got %d", or_.Len())
	}
}
