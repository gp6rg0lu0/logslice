package filter_test

import (
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeWildcardLogLine(field, value string) *parser.LogLine {
	data := map[string]string{field: value}
	return parser.NewLogLine(data)
}

func TestNewWildcardFilter_EmptyField(t *testing.T) {
	_, err := filter.NewWildcardFilter("", "foo*")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewWildcardFilter_EmptyPattern(t *testing.T) {
	_, err := filter.NewWildcardFilter("msg", "")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNewWildcardFilter_InvalidPattern(t *testing.T) {
	_, err := filter.NewWildcardFilter("msg", "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid glob pattern")
	}
}

func TestWildcardFilter_Accessors(t *testing.T) {
	f, err := filter.NewWildcardFilter("msg", "err*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "msg" {
		t.Errorf("expected field 'msg', got %q", f.Field())
	}
	if f.Pattern() != "err*" {
		t.Errorf("expected pattern 'err*', got %q", f.Pattern())
	}
}

func TestWildcardFilter_NilLine(t *testing.T) {
	f, _ := filter.NewWildcardFilter("msg", "err*")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestWildcardFilter_MissingField(t *testing.T) {
	f, _ := filter.NewWildcardFilter("msg", "err*")
	line := makeWildcardLogLine("other", "error occurred")
	if f.Match(line) {
		t.Error("expected false for missing field")
	}
}

func TestWildcardFilter_Match(t *testing.T) {
	tests := []struct {
		pattern string
		value   string
		want    bool
	}{
		{"err*", "error occurred", true},
		{"err*", "warning", false},
		{"*fail*", "connection failed", true},
		{"exact", "exact", true},
		{"exact", "not exact", false},
		{"*.json", "config.json", true},
		{"*.json", "config.yaml", false},
	}
	for _, tt := range tests {
		f, err := filter.NewWildcardFilter("msg", tt.pattern)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		line := makeWildcardLogLine("msg", tt.value)
		if got := f.Match(line); got != tt.want {
			t.Errorf("pattern=%q value=%q: expected %v, got %v", tt.pattern, tt.value, tt.want, got)
		}
	}
}
