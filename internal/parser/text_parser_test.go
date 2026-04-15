package parser

import (
	"io"
	"strings"
	"testing"
)

const samplePattern = `(?P<time>\d{4}-\d{2}-\d{2}T\S+) (?P<level>\w+) (?P<message>.*)`

func TestNewTextParser_InvalidPattern(t *testing.T) {
	_, err := NewTextParser(strings.NewReader(""), "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNewTextParser_NoNamedGroups(t *testing.T) {
	_, err := NewTextParser(strings.NewReader(""), `\w+`)
	if err == nil {
		t.Fatal("expected error when pattern has no named groups")
	}
}

func TestTextParser_SingleLine(t *testing.T) {
	input := "2024-01-15T10:00:00Z INFO server started"
	p, err := NewTextParser(strings.NewReader(input), samplePattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line, err := p.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := line.Get("level"); got != "INFO" {
		t.Errorf("level: got %q, want %q", got, "INFO")
	}
	if got := line.Get("message"); got != "server started" {
		t.Errorf("message: got %q, want %q", got, "server started")
	}
	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestTextParser_SkipsNonMatchingLines(t *testing.T) {
	input := "not a log line\n2024-01-15T10:00:00Z WARN disk full\nalso not matching"
	p, err := NewTextParser(strings.NewReader(input), samplePattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line, err := p.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := line.Get("level"); got != "WARN" {
		t.Errorf("level: got %q, want %q", got, "WARN")
	}
	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected EOF after single match, got %v", err)
	}
}

func TestTextParser_SkipsBlankLines(t *testing.T) {
	input := "\n\n2024-01-15T10:00:00Z ERROR boom\n\n"
	p, err := NewTextParser(strings.NewReader(input), samplePattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line, err := p.Next()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := line.Get("level"); got != "ERROR" {
		t.Errorf("level: got %q, want %q", got, "ERROR")
	}
}

func TestTextParser_MultipleLines(t *testing.T) {
	input := "2024-01-15T10:00:00Z DEBUG a\n2024-01-15T10:00:01Z INFO b\n2024-01-15T10:00:02Z ERROR c"
	p, err := NewTextParser(strings.NewReader(input), samplePattern)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	levels := []string{"DEBUG", "INFO", "ERROR"}
	for _, want := range levels {
		line, err := p.Next()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := line.Get("level"); got != want {
			t.Errorf("level: got %q, want %q", got, want)
		}
	}
}
