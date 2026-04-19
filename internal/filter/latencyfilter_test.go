package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeLatencyLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]string{field: value})
}

func TestNewLatencyFilter_EmptyField(t *testing.T) {
	_, err := NewLatencyFilter("", "0:100")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewLatencyFilter_BadFormat(t *testing.T) {
	_, err := NewLatencyFilter("latency_ms", "100")
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestNewLatencyFilter_InvalidMin(t *testing.T) {
	_, err := NewLatencyFilter("latency_ms", "abc:100")
	if err == nil {
		t.Fatal("expected error for invalid min")
	}
}

func TestNewLatencyFilter_InvalidMax(t *testing.T) {
	_, err := NewLatencyFilter("latency_ms", "0:xyz")
	if err == nil {
		t.Fatal("expected error for invalid max")
	}
}

func TestNewLatencyFilter_MinExceedsMax(t *testing.T) {
	_, err := NewLatencyFilter("latency_ms", "200:100")
	if err == nil {
		t.Fatal("expected error when min > max")
	}
}

func TestLatencyFilter_Accessors(t *testing.T) {
	f, err := NewLatencyFilter("latency_ms", "10:500")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "latency_ms" {
		t.Errorf("expected field latency_ms, got %s", f.Field())
	}
	if f.Min() != 10 || f.Max() != 500 {
		t.Errorf("unexpected bounds: min=%.0f max=%.0f", f.Min(), f.Max())
	}
}

func TestLatencyFilter_NilLine(t *testing.T) {
	f, _ := NewLatencyFilter("latency_ms", "0:100")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestLatencyFilter_Match(t *testing.T) {
	f, _ := NewLatencyFilter("latency_ms", "50:200")
	tests := []struct {
		val  string
		want bool
	}{
		{"100", true},
		{"50", true},
		{"200", true},
		{"49", false},
		{"201", false},
		{"notanumber", false},
	}
	for _, tc := range tests {
		line := makeLatencyLogLine("latency_ms", tc.val)
		if got := f.Match(line); got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.val, got, tc.want)
		}
	}
}

func TestLatencyFilter_Unbounded(t *testing.T) {
	f, err := NewLatencyFilter("latency_ms", ":")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeLatencyLogLine("latency_ms", "9999")
	if !f.Match(line) {
		t.Error("expected match for unbounded filter")
	}
}
