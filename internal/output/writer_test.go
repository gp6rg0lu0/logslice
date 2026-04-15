package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
)

func makeLogLine(t *testing.T, raw string) *parser.LogLine {
	t.Helper()
	p := parser.NewJSONParser(strings.NewReader(raw))
	line, err := p.Next()
	if err != nil {
		t.Fatalf("makeLogLine: %v", err)
	}
	return line
}

func TestNewWriter_DefaultsToJSON(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.Format("unknown"))
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestWriter_WriteNil(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatJSON)
	if err := w.Write(nil); err != nil {
		t.Fatalf("expected no error writing nil, got %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected empty output for nil line")
	}
}

func TestWriter_WriteJSON(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatJSON)
	line := makeLogLine(t, `{"level":"info","msg":"hello"}`)
	if err := w.Write(line); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result["level"] != "info" {
		t.Errorf("expected level=info, got %v", result["level"])
	}
}

func TestWriter_WriteText(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatText)
	line := makeLogLine(t, `{"level":"warn","msg":"oops"}`)
	if err := w.Write(line); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "level=warn") {
		t.Errorf("expected 'level=warn' in text output, got: %s", out)
	}
}

func TestWriter_WriteRaw(t *testing.T) {
	const raw = `{"level":"debug","msg":"trace"}`
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatRaw)
	line := makeLogLine(t, raw)
	if err := w.Write(line); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), raw) {
		t.Errorf("expected raw line in output, got: %s", buf.String())
	}
}
