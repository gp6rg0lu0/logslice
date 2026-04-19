package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func makeHTTPStatusLogLine(field, value string) *parser.LogLine {
	if field == "" {
		return parser.NewLogLine(map[string]interface{}{})
	}
	return parser.NewLogLine(map[string]interface{}{field: value})
}

func TestNewHTTPStatusFilter_EmptyField(t *testing.T) {
	_, err := filter.NewHTTPStatusFilter("", 200, 299)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewHTTPStatusFilter_InvalidRange(t *testing.T) {
	if _, err := filter.NewHTTPStatusFilter("status", 50, 200); err == nil {
		t.Fatal("expected error for min < 100")
	}
	if _, err := filter.NewHTTPStatusFilter("status", 200, 600); err == nil {
		t.Fatal("expected error for max > 599")
	}
	if _, err := filter.NewHTTPStatusFilter("status", 400, 300); err == nil {
		t.Fatal("expected error for min > max")
	}
}

func TestHTTPStatusFilter_Accessors(t *testing.T) {
	f, err := filter.NewHTTPStatusFilter("status", 400, 499)
	if err != nil {
		t.Fatal(err)
	}
	if f.Field() != "status" {
		t.Errorf("expected 'status', got %q", f.Field())
	}
	if f.Min() != 400 || f.Max() != 499 {
		t.Errorf("unexpected bounds: %d %d", f.Min(), f.Max())
	}
}

func TestHTTPStatusFilter_NilLine(t *testing.T) {
	f, _ := filter.NewHTTPStatusFilter("status", 200, 299)
	if f.Match(nil) {
		t.Fatal("expected false for nil line")
	}
}

func TestHTTPStatusFilter_Match(t *testing.T) {
	f, _ := filter.NewHTTPStatusFilter("status", 400, 499)
	cases := []struct {
		val  string
		want bool
	}{
		{"400", true},
		{"404", true},
		{"499", true},
		{"200", false},
		{"500", false},
		{"abc", false},
		{"", false},
	}
	for _, c := range cases {
		line := makeHTTPStatusLogLine("status", c.val)
		if got := f.Match(line); got != c.want {
			t.Errorf("Match(%q) = %v, want %v", c.val, got, c.want)
		}
	}
}
