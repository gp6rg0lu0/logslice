package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEnvLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]interface{}{field: value})
}

func TestNewEnvFilter_EmptyField(t *testing.T) {
	_, err := NewEnvFilter("", []string{"prod"}, false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewEnvFilter_NoValues(t *testing.T) {
	_, err := NewEnvFilter("env", []string{}, false)
	if err == nil {
		t.Fatal("expected error for empty envs")
	}
}

func TestNewEnvFilter_BlankValue(t *testing.T) {
	_, err := NewEnvFilter("env", []string{"prod", "  "}, false)
	if err == nil {
		t.Fatal("expected error for blank env value")
	}
}

func TestEnvFilter_Accessors(t *testing.T) {
	f, err := NewEnvFilter("env", []string{"prod", "staging"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "env" {
		t.Errorf("expected field 'env', got %q", f.Field())
	}
	if !f.CaseFold() {
		t.Error("expected caseFold true")
	}
	if len(f.Envs()) != 2 {
		t.Errorf("expected 2 envs, got %d", len(f.Envs()))
	}
}

func TestEnvFilter_NilLine(t *testing.T) {
	f, _ := NewEnvFilter("env", []string{"prod"}, false)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestEnvFilter_MissingField(t *testing.T) {
	f, _ := NewEnvFilter("env", []string{"prod"}, false)
	line := makeEnvLogLine("other", "prod")
	if f.Match(line) {
		t.Error("expected false when field missing")
	}
}

func TestEnvFilter_MatchExact(t *testing.T) {
	f, _ := NewEnvFilter("env", []string{"prod", "staging"}, false)
	if !f.Match(makeEnvLogLine("env", "prod")) {
		t.Error("expected match for 'prod'")
	}
	if !f.Match(makeEnvLogLine("env", "staging")) {
		t.Error("expected match for 'staging'")
	}
	if f.Match(makeEnvLogLine("env", "dev")) {
		t.Error("expected no match for 'dev'")
	}
}

func TestEnvFilter_CaseFold(t *testing.T) {
	f, _ := NewEnvFilter("env", []string{"PROD"}, true)
	if !f.Match(makeEnvLogLine("env", "prod")) {
		t.Error("expected case-insensitive match")
	}
	if !f.Match(makeEnvLogLine("env", "Prod")) {
		t.Error("expected case-insensitive match for mixed case")
	}
}
