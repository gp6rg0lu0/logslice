package filter

import (
	"testing"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		input    string
		expected LogLevel
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"INFO", LevelInfo},
		{"warn", LevelWarn},
		{"warning", LevelWarn},
		{"error", LevelError},
		{"err", LevelError},
		{"fatal", LevelFatal},
		{"panic", LevelFatal},
		{"unknown", LevelUnknown},
		{"", LevelUnknown},
	}
	for _, tc := range cases {
		got := ParseLevel(tc.input)
		if got != tc.expected {
			t.Errorf("ParseLevel(%q) = %d, want %d", tc.input, got, tc.expected)
		}
	}
}

func TestLevelFilter_Match(t *testing.T) {
	cases := []struct {
		minLevel string
		line     string
		expected bool
	}{
		{"info", `{"level":"info","msg":"started"}`, true},
		{"info", `{"level":"debug","msg":"verbose"}`, false},
		{"info", `{"level":"error","msg":"oops"}`, true},
		{"warn", `{"level":"warn","msg":"slow"}`, true},
		{"warn", `{"level":"info","msg":"ok"}`, false},
		{"error", `{"lvl":"error","msg":"fail"}`, true},
		{"error", `{"severity":"fatal","msg":"crash"}`, true},
		{"debug", `{"level":"debug","msg":"trace"}`, true},
		{"info", `not json at all`, false},
		{"info", `{"msg":"no level field"}`, false},
	}
	for _, tc := range cases {
		f := NewLevelFilter(tc.minLevel)
		got := f.Match(tc.line)
		if got != tc.expected {
			t.Errorf("LevelFilter(%q).Match(%q) = %v, want %v",
				tc.minLevel, tc.line, got, tc.expected)
		}
	}
}
