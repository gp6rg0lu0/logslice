package filter

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makeSeverityLogLine(field, value string) *parser.LogLine {
	if field == "" {
		return parser.NewLogLine(map[string]string{})
	}
	return parser.NewLogLine(map[string]string{field: value})
}

func TestNewSeverityFilter_EmptyField(t *testing.T) {
	_, err := NewSeverityFilter("", "debug", "error")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewSeverityFilter_UnknownMin(t *testing.T) {
	_, err := NewSeverityFilter("severity", "verbose", "error")
	if err == nil {
		t.Fatal("expected error for unknown min severity")
	}
}

func TestNewSeverityFilter_UnknownMax(t *testing.T) {
	_, err := NewSeverityFilter("severity", "debug", "catastrophic")
	if err == nil {
		t.Fatal("expected error for unknown max severity")
	}
}

func TestNewSeverityFilter_MinExceedsMax(t *testing.T) {
	_, err := NewSeverityFilter("severity", "error", "debug")
	if err == nil {
		t.Fatal("expected error when min exceeds max")
	}
}

func TestNewSeverityFilter_Defaults(t *testing.T) {
	f, err := NewSeverityFilter("sev", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Min() != "trace" {
		t.Errorf("expected min=trace, got %q", f.Min())
	}
	if f.Max() != "panic" {
		t.Errorf("expected max=panic, got %q", f.Max())
	}
}

func TestSeverityFilter_Accessors(t *testing.T) {
	f, err := NewSeverityFilter("severity", "info", "error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "severity" {
		t.Errorf("expected field=severity, got %q", f.Field())
	}
	if f.Min() != "info" {
		t.Errorf("expected min=info, got %q", f.Min())
	}
	if f.Max() != "error" {
		t.Errorf("expected max=error, got %q", f.Max())
	}
}

func TestSeverityFilter_NilLine(t *testing.T) {
	f, _ := NewSeverityFilter("severity", "debug", "error")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestSeverityFilter_MissingField(t *testing.T) {
	f, _ := NewSeverityFilter("severity", "debug", "error")
	line := makeSeverityLogLine("", "")
	if f.Match(line) {
		t.Error("expected false for missing field")
	}
}

func TestSeverityFilter_UnknownValue(t *testing.T) {
	f, _ := NewSeverityFilter("severity", "debug", "error")
	line := makeSeverityLogLine("severity", "verbose")
	if f.Match(line) {
		t.Error("expected false for unknown severity value")
	}
}

func TestSeverityFilter_Match(t *testing.T) {
	f, _ := NewSeverityFilter("severity", "warning", "critical")
	cases := []struct {
		value string
		want  bool
	}{
		{"trace", false},
		{"debug", false},
		{"info", false},
		{"warning", true},
		{"WARN", true},
		{"error", true},
		{"critical", true},
		{"fatal", false},
		{"panic", false},
	}
	for _, tc := range cases {
		line := makeSeverityLogLine("severity", tc.value)
		if got := f.Match(line); got != tc.want {
			t.Errorf("value=%q: got %v, want %v", tc.value, got, tc.want)
		}
	}
}
