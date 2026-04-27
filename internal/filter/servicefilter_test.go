package filter

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func makeServiceLogLine(field, value string) *parser.LogLine {
	fields := map[string]interface{}{}
	if field != "" {
		fields[field] = value
	}
	return parser.NewLogLine(fields)
}

func TestNewServiceFilter_EmptyField(t *testing.T) {
	_, err := NewServiceFilter("", []string{"api"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewServiceFilter_NoServices(t *testing.T) {
	_, err := NewServiceFilter("service", []string{})
	if err == nil {
		t.Fatal("expected error for empty services slice")
	}
}

func TestNewServiceFilter_BlankService(t *testing.T) {
	_, err := NewServiceFilter("service", []string{"api", "  "})
	if err == nil {
		t.Fatal("expected error for blank service name")
	}
}

func TestServiceFilter_Accessors(t *testing.T) {
	f, err := NewServiceFilter("svc", []string{"auth", "billing"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "svc" {
		t.Errorf("expected field 'svc', got %q", f.Field())
	}
	if len(f.Services()) != 2 {
		t.Errorf("expected 2 services, got %d", len(f.Services()))
	}
}

func TestServiceFilter_NilLine(t *testing.T) {
	f, _ := NewServiceFilter("service", []string{"api"})
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestServiceFilter_MatchesExact(t *testing.T) {
	f, _ := NewServiceFilter("service", []string{"api", "worker"})
	line := makeServiceLogLine("service", "api")
	if !f.Match(line) {
		t.Error("expected match for 'api'")
	}
}

func TestServiceFilter_MatchesCaseInsensitive(t *testing.T) {
	f, _ := NewServiceFilter("service", []string{"API"})
	line := makeServiceLogLine("service", "api")
	if !f.Match(line) {
		t.Error("expected case-insensitive match")
	}
}

func TestServiceFilter_NoMatch(t *testing.T) {
	f, _ := NewServiceFilter("service", []string{"auth"})
	line := makeServiceLogLine("service", "billing")
	if f.Match(line) {
		t.Error("expected no match for 'billing'")
	}
}

func TestServiceFilter_MissingField(t *testing.T) {
	f, _ := NewServiceFilter("service", []string{"api"})
	line := makeServiceLogLine("", "")
	if f.Match(line) {
		t.Error("expected no match when field is absent")
	}
}
