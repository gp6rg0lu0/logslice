package filter

import (
	"testing"

	"github.com/nicholasgasior/logslice/internal/parser"
)

func makeJSONPathLogLine(fields map[string]interface{}) *parser.LogLine {
	return parser.NewLogLineFromFields(fields)
}

func TestNewJSONPathFilter_EmptyPath(t *testing.T) {
	_, err := NewJSONPathFilter("", "val")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNewJSONPathFilter_EmptyValue(t *testing.T) {
	_, err := NewJSONPathFilter("a.b", "")
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestNewJSONPathFilter_EmptySegment(t *testing.T) {
	_, err := NewJSONPathFilter("a..b", "val")
	if err == nil {
		t.Fatal("expected error for empty segment")
	}
}

func TestJSONPathFilter_Accessors(t *testing.T) {
	f, err := NewJSONPathFilter("a.b", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Path() != "a.b" {
		t.Errorf("expected path a.b, got %s", f.Path())
	}
	if f.Value() != "hello" {
		t.Errorf("expected value hello, got %s", f.Value())
	}
}

func TestJSONPathFilter_NilLine(t *testing.T) {
	f, _ := NewJSONPathFilter("a.b", "v")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestJSONPathFilter_FlatMatch(t *testing.T) {
	f, _ := NewJSONPathFilter("env", "production")
	line := makeJSONPathLogLine(map[string]interface{}{"env": "production"})
	if !f.Match(line) {
		t.Error("expected match")
	}
}

func TestJSONPathFilter_NestedMatch(t *testing.T) {
	f, _ := NewJSONPathFilter("meta.region", "us-east-1")
	line := makeJSONPathLogLine(map[string]interface{}{
		"meta": map[string]interface{}{"region": "us-east-1"},
	})
	if !f.Match(line) {
		t.Error("expected nested match")
	}
}

func TestJSONPathFilter_NestedMismatch(t *testing.T) {
	f, _ := NewJSONPathFilter("meta.region", "eu-west-1")
	line := makeJSONPathLogLine(map[string]interface{}{
		"meta": map[string]interface{}{"region": "us-east-1"},
	})
	if f.Match(line) {
		t.Error("expected no match")
	}
}

func TestJSONPathFilter_MissingKey(t *testing.T) {
	f, _ := NewJSONPathFilter("meta.region", "us-east-1")
	line := makeJSONPathLogLine(map[string]interface{}{"env": "prod"})
	if f.Match(line) {
		t.Error("expected no match for missing key")
	}
}
