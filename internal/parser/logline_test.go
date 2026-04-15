package parser

import "testing"

func TestLogLine_Get_NilReceiver(t *testing.T) {
	var l *LogLine
	if got := l.Get("level"); got != "" {
		t.Fatalf("expected empty string from nil receiver, got %q", got)
	}
}

func TestLogLine_Get_MissingKey(t *testing.T) {
	l := &LogLine{Fields: map[string]string{"level": "info"}}
	if got := l.Get("missing"); got != "" {
		t.Fatalf("expected empty string for missing key, got %q", got)
	}
}

func TestLogLine_Get_PresentKey(t *testing.T) {
	l := &LogLine{Fields: map[string]string{"msg": "hello"}}
	if got := l.Get("msg"); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}

func TestLogLine_Level(t *testing.T) {
	l := &LogLine{Fields: map[string]string{"level": "warn"}}
	if got := l.Level(); got != "warn" {
		t.Fatalf("expected %q, got %q", "warn", got)
	}
}

func TestLogLine_Time(t *testing.T) {
	ts := "2024-01-15T10:00:00Z"
	l := &LogLine{Fields: map[string]string{"time": ts}}
	if got := l.Time(); got != ts {
		t.Fatalf("expected %q, got %q", ts, got)
	}
}

func TestLogLine_NilFields(t *testing.T) {
	l := &LogLine{Raw: "raw line"}
	if got := l.Level(); got != "" {
		t.Fatalf("expected empty string when Fields is nil, got %q", got)
	}
	if got := l.Time(); got != "" {
		t.Fatalf("expected empty string when Fields is nil, got %q", got)
	}
}

func TestLogLine_RawPreserved(t *testing.T) {
	raw := `{"level":"error","msg":"boom"}`
	l := &LogLine{
		Raw:    raw,
		Fields: map[string]string{"level": "error", "msg": "boom"},
	}
	if l.Raw != raw {
		t.Fatalf("raw field not preserved: got %q", l.Raw)
	}
}
