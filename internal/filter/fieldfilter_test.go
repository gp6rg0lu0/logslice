package filter

import (
	"testing"
)

func TestNewFieldFilter_InvalidArgs(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		pattern string
	}{
		{"empty field", "", "foo"},
		{"empty pattern", "level", ""},
		{"invalid regex", "msg", "[unclosed"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewFieldFilter(tc.field, tc.pattern)
			if err == nil {
				t.Errorf("expected error for field=%q pattern=%q, got nil", tc.field, tc.pattern)
			}
		})
	}
}

func TestFieldFilter_Match(t *testing.T) {
	f, err := NewFieldFilter("service", "^auth.*")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		name  string
		entry map[string]interface{}
		want  bool
	}{
		{"matching string value", map[string]interface{}{"service": "auth-api"}, true},
		{"non-matching string value", map[string]interface{}{"service": "gateway"}, false},
		{"missing field", map[string]interface{}{"level": "info"}, false},
		{"numeric value matching pattern", map[string]interface{}{"service": 42}, false},
		{"string auth prefix exact", map[string]interface{}{"service": "auth"}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := f.Match(tc.entry)
			if got != tc.want {
				t.Errorf("Match(%v) = %v, want %v", tc.entry, got, tc.want)
			}
		})
	}
}

func TestFieldFilter_Accessors(t *testing.T) {
	f, err := NewFieldFilter("msg", "error.*timeout")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "msg" {
		t.Errorf("Field() = %q, want %q", f.Field(), "msg")
	}
	if f.Pattern() != "error.*timeout" {
		t.Errorf("Pattern() = %q, want %q", f.Pattern(), "error.*timeout")
	}
}
