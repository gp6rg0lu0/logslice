package parser_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func TestNew_JSONFormat(t *testing.T) {
	p, err := parser.New(parser.FormatJSON, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil parser")
	}
}

func TestNew_TextFormat(t *testing.T) {
	pattern := `(?P<time>\S+) (?P<level>\S+) (?P<message>.+)`
	p, err := parser.New(parser.FormatText, pattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil parser")
	}
}

func TestNew_TextFormat_InvalidPattern(t *testing.T) {
	_, err := parser.New(parser.FormatText, "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid pattern, got nil")
	}
}

func TestNew_UnknownFormat(t *testing.T) {
	_, err := parser.New("xml", "")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
	var ufErr *parser.UnknownFormatError
	// Check error message contains format name
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("expected error to mention format name, got: %v", err)
	}
	_ = ufErr
}

func TestNew_JSONFormat_ParsesLines(t *testing.T) {
	p, err := parser.New(parser.FormatJSON, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := `{"level":"info","msg":"hello"}
{"level":"error","msg":"world"}
`
	r := strings.NewReader(input)
	ch := p.Parse(r)

	var lines []*parser.LogLine
	for line := range ch {
		lines = append(lines, line)
	}

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if v, _ := lines[0].Get("msg"); v != "hello" {
		t.Errorf("expected msg=hello, got %q", v)
	}
}

func TestUnknownFormatError_Error(t *testing.T) {
	err := &parser.UnknownFormatError{Format: "csv"}
	if !strings.Contains(err.Error(), "csv") {
		t.Errorf("expected error message to contain 'csv', got: %s", err.Error())
	}
}
