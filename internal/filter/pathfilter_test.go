package filter

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makePathLogLine(field, value string) *parser.LogLine {
	fields := map[string]interface{}{field: value}
	return parser.NewLogLine(fields)
}

func TestNewPathFilter_EmptyField(t *testing.T) {
	_, err := NewPathFilter("", []string{"/api"}, false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewPathFilter_NoPaths(t *testing.T) {
	_, err := NewPathFilter("path", []string{}, false)
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestNewPathFilter_BlankPath(t *testing.T) {
	_, err := NewPathFilter("path", []string{"  "}, false)
	if err == nil {
		t.Fatal("expected error for blank path")
	}
}

func TestPathFilter_Accessors(t *testing.T) {
	f, err := NewPathFilter("url", []string{"/api", "/v2"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "url" {
		t.Errorf("expected field 'url', got %q", f.Field())
	}
	if len(f.Paths()) != 2 {
		t.Errorf("expected 2 paths, got %d", len(f.Paths()))
	}
	if !f.Exact() {
		t.Error("expected exact=true")
	}
}

func TestPathFilter_NilLine(t *testing.T) {
	f, _ := NewPathFilter("path", []string{"/api"}, false)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestPathFilter_MissingField(t *testing.T) {
	f, _ := NewPathFilter("path", []string{"/api"}, false)
	line := makePathLogLine("url", "/api/users")
	if f.Match(line) {
		t.Error("expected false when field is missing")
	}
}

func TestPathFilter_PrefixMatch(t *testing.T) {
	f, _ := NewPathFilter("path", []string{"/api/v1", "/health"}, false)
	cases := []struct {
		val  string
		want bool
	}{
		{"/api/v1/users", true},
		{"/health", true},
		{"/api/v2/items", false},
		{"/metrics", false},
	}
	for _, tc := range cases {
		line := makePathLogLine("path", tc.val)
		if got := f.Match(line); got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.val, got, tc.want)
		}
	}
}

func TestPathFilter_ExactMatch(t *testing.T) {
	f, _ := NewPathFilter("path", []string{"/health", "/ready"}, true)
	cases := []struct {
		val  string
		want bool
	}{
		{"/health", true},
		{"/health/check", false},
		{"/ready", true},
		{"/readyz", false},
	}
	for _, tc := range cases {
		line := makePathLogLine("path", tc.val)
		if got := f.Match(line); got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.val, got, tc.want)
		}
	}
}
