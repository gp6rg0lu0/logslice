package output

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeFormatterLogLine(fields map[string]interface{}) *parser.LogLine {
	raw, _ := json.Marshal(fields)
	p := parser.NewJSONParser()
	lines, _ := p.Parse(strings.NewReader(string(raw)))
	if len(lines) == 0 {
		return nil
	}
	return lines[0]
}

func TestJSONFormatter_NilLine(t *testing.T) {
	f := NewJSONFormatter()
	_, err := f.Format(nil)
	if err == nil {
		t.Fatal("expected error for nil line")
	}
}

func TestJSONFormatter_Format(t *testing.T) {
	f := NewJSONFormatter()
	line := makeFormatterLogLine(map[string]interface{}{"level": "info", "msg": "hello"})
	if line == nil {
		t.Fatal("failed to create log line")
	}
	out, err := f.Format(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
}

func TestTextFormatter_NilLine(t *testing.T) {
	f := NewTextFormatter(nil)
	_, err := f.Format(nil)
	if err == nil {
		t.Fatal("expected error for nil line")
	}
}

func TestTextFormatter_PriorityKeys(t *testing.T) {
	f := NewTextFormatter([]string{"time", "level"})
	line := makeFormatterLogLine(map[string]interface{}{"level": "warn", "time": "2024-01-01T00:00:00Z", "msg": "test"})
	if line == nil {
		t.Fatal("failed to create log line")
	}
	out, err := f.Format(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "time=") {
		t.Errorf("expected output to start with time=, got: %s", out)
	}
	if !strings.Contains(out, "level=warn") {
		t.Errorf("expected level=warn in output: %s", out)
	}
}

func TestTextFormatter_AlphabeticRemainder(t *testing.T) {
	f := NewTextFormatter([]string{})
	line := makeFormatterLogLine(map[string]interface{}{"zebra": "z", "apple": "a"})
	if line == nil {
		t.Fatal("failed to create log line")
	}
	out, err := f.Format(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	appleIdx := strings.Index(out, "apple=")
	zebraIdx := strings.Index(out, "zebra=")
	if appleIdx > zebraIdx {
		t.Errorf("expected apple before zebra in: %s", out)
	}
}
