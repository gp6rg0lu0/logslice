package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func makeNotLogLine(fields map[string]string) *parser.LogLine {
	raw := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		raw[k] = v
	}
	return parser.NewLogLine(raw)
}

func TestNewNotFilter_NilInner(t *testing.T) {
	if filter.NewNotFilter(nil) != nil {
		t.Fatal("expected nil for nil inner")
	}
}

func TestNotFilter_InvertsMatch(t *testing.T) {
	inner, err := filter.NewLevelFilter("info")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	not := filter.NewNotFilter(inner)

	tests := []struct {
		level string
		want  bool
	}{
		{"info", false},
		{"debug", true},
		{"error", true},
	}
	for _, tt := range tests {
		line := makeNotLogLine(map[string]string{"level": tt.level})
		if got := not.Match(line); got != tt.want {
			t.Errorf("level=%q: got %v, want %v", tt.level, got, tt.want)
		}
	}
}

func TestNotFilter_Inner(t *testing.T) {
	inner, _ := filter.NewLevelFilter("warn")
	not := filter.NewNotFilter(inner)
	if not.Inner() != inner {
		t.Fatal("Inner() did not return wrapped filter")
	}
}

func TestNotFilter_WithChain(t *testing.T) {
	inner, _ := filter.NewLevelFilter("error")
	not := filter.NewNotFilter(inner)
	chain := filter.NewChain(not)

	line := makeNotLogLine(map[string]string{"level": "info"})
	if !chain.Match(line) {
		t.Error("chain with NotFilter should match non-error line")
	}

	errLine := makeNotLogLine(map[string]string{"level": "error"})
	if chain.Match(errLine) {
		t.Error("chain with NotFilter should not match error line")
	}
}
