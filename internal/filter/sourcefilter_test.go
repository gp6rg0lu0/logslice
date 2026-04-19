package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeSourceLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]interface{}{field: value})
}

func TestNewSourceFilter_EmptyField(t *testing.T) {
	_, err := NewSourceFilter("", []string{"app"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewSourceFilter_NoValues(t *testing.T) {
	_, err := NewSourceFilter("source", []string{})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}

func TestNewSourceFilter_BlankValues(t *testing.T) {
	_, err := NewSourceFilter("source", []string{"  ", ""})
	if err == nil {
		t.Fatal("expected error when all sources are blank")
	}
}

func TestSourceFilter_Accessors(t *testing.T) {
	f, err := NewSourceFilter("src", []string{"app", "worker"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "src" {
		t.Errorf("expected field 'src', got %q", f.Field())
	}
	if len(f.Sources()) != 2 {
		t.Errorf("expected 2 sources, got %d", len(f.Sources()))
	}
}

func TestSourceFilter_NilLine(t *testing.T) {
	f, _ := NewSourceFilter("source", []string{"app"})
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestSourceFilter_MatchesCaseInsensitive(t *testing.T) {
	f, _ := NewSourceFilter("source", []string{"App", "Worker"})
	if !f.Match(makeSourceLogLine("source", "app")) {
		t.Error("expected match for 'app'")
	}
	if !f.Match(makeSourceLogLine("source", "WORKER")) {
		t.Error("expected match for 'WORKER'")
	}
}

func TestSourceFilter_NoMatch(t *testing.T) {
	f, _ := NewSourceFilter("source", []string{"app"})
	if f.Match(makeSourceLogLine("source", "db")) {
		t.Error("expected no match for 'db'")
	}
}

func TestSourceFilter_MissingField(t *testing.T) {
	f, _ := NewSourceFilter("source", []string{"app"})
	line := parser.NewLogLine(map[string]interface{}{"other": "app"})
	if f.Match(line) {
		t.Error("expected no match when field is missing")
	}
}
