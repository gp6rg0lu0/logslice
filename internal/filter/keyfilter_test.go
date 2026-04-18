package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeKeyLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]interface{}{field: value})
}

func TestNewKeyFilter_EmptyField(t *testing.T) {
	_, err := NewKeyFilter("", "a,b")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewKeyFilter_EmptyKeys(t *testing.T) {
	_, err := NewKeyFilter("env", "")
	if err == nil {
		t.Fatal("expected error for empty keys")
	}
}

func TestNewKeyFilter_BlankKeys(t *testing.T) {
	_, err := NewKeyFilter("env", "  ,  ")
	if err == nil {
		t.Fatal("expected error when all keys are blank")
	}
}

func TestKeyFilter_Accessors(t *testing.T) {
	f, err := NewKeyFilter("env", "prod,staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "env" {
		t.Errorf("expected field 'env', got %q", f.Field())
	}
	if f.Keys() != "prod,staging" {
		t.Errorf("expected keys 'prod,staging', got %q", f.Keys())
	}
}

func TestKeyFilter_NilLine(t *testing.T) {
	f, _ := NewKeyFilter("env", "prod")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestKeyFilter_MatchesExact(t *testing.T) {
	f, _ := NewKeyFilter("env", "prod,staging")
	if !f.Match(makeKeyLogLine("env", "prod")) {
		t.Error("expected match for 'prod'")
	}
	if !f.Match(makeKeyLogLine("env", "staging")) {
		t.Error("expected match for 'staging'")
	}
}

func TestKeyFilter_NoMatchOtherValue(t *testing.T) {
	f, _ := NewKeyFilter("env", "prod,staging")
	if f.Match(makeKeyLogLine("env", "dev")) {
		t.Error("expected no match for 'dev'")
	}
}

func TestKeyFilter_MissingField(t *testing.T) {
	f, _ := NewKeyFilter("env", "prod")
	line := parser.NewLogLine(map[string]interface{}{"other": "prod"})
	if f.Match(line) {
		t.Error("expected no match when field is missing")
	}
}
