package parser

import (
	"io"
	"strings"
	"testing"
)

func TestJSONParser_SingleLine(t *testing.T) {
	input := `{"level":"info","msg":"started","ts":"2024-01-01T00:00:00Z"}`
	p := NewJSONParser(strings.NewReader(input))
	entry, err := p.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry["level"] != "info" {
		t.Errorf("expected level=info, got %q", entry["level"])
	}
	if entry["msg"] != "started" {
		t.Errorf("expected msg=started, got %q", entry["msg"])
	}
}

func TestJSONParser_MultipleLines(t *testing.T) {
	input := "{\"level\":\"info\"}\n{\"level\":\"error\"}\n"
	p := NewJSONParser(strings.NewReader(input))
	levels := []string{"info", "error"}
	for _, want := range levels {
		entry, err := p.Next()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if entry["level"] != want {
			t.Errorf("expected %q, got %q", want, entry["level"])
		}
	}
	_, err := p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestJSONParser_SkipsBlankLines(t *testing.T) {
	input := "\n{\"level\":\"warn\"}\n\n"
	p := NewJSONParser(strings.NewReader(input))
	entry, err := p.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry["level"] != "warn" {
		t.Errorf("expected warn, got %q", entry["level"])
	}
}

func TestJSONParser_NonStringValues(t *testing.T) {
	input := `{"code":404,"ok":false}`
	p := NewJSONParser(strings.NewReader(input))
	entry, err := p.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry["code"] != "404" {
		t.Errorf("expected \"404\", got %q", entry["code"])
	}
	if entry["ok"] != "false" {
		t.Errorf("expected \"false\", got %q", entry["ok"])
	}
}

func TestJSONParser_InvalidJSON(t *testing.T) {
	input := "not json\n"
	p := NewJSONParser(strings.NewReader(input))
	_, err := p.Next()
	if err == nil {
		t.Error("expected parse error for invalid JSON")
	}
}

func TestJSONParser_Empty(t *testing.T) {
	p := NewJSONParser(strings.NewReader(""))
	_, err := p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF on empty input, got %v", err)
	}
}
