package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeRegexLogLine(fields map[string]string) *parser.LogLine {
	raw := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		raw[k] = v
	}
	return parser.NewLogLine(raw)
}

func TestNewRegexFilter_EmptyField(t *testing.T) {
	_, err := NewRegexFilter("", "foo")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewRegexFilter_EmptyPattern(t *testing.T) {
	_, err := NewRegexFilter("msg", "")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNewRegexFilter_InvalidPattern(t *testing.T) {
	_, err := NewRegexFilter("msg", "[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestRegexFilter_Accessors(t *testing.T) {
	f, err := NewRegexFilter("msg", `^error`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "msg" {
		t.Errorf("expected field 'msg', got %q", f.Field())
	}
	if f.Pattern().String() != `^error` {
		t.Errorf("unexpected pattern: %s", f.Pattern())
	}
}

func TestRegexFilter_Match_NilLine(t *testing.T) {
	f, _ := NewRegexFilter("msg", `error`)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestRegexFilter_Match_MissingField(t *testing.T) {
	f, _ := NewRegexFilter("msg", `error`)
	line := makeRegexLogLine(map[string]string{"level": "info"})
	if f.Match(line) {
		t.Error("expected false when field is absent")
	}
}

func TestRegexFilter_Match_Matches(t *testing.T) {
	f, _ := NewRegexFilter("msg", `(?i)timeout`)
	line := makeRegexLogLine(map[string]string{"msg": "connection Timeout exceeded"})
	if !f.Match(line) {
		t.Error("expected true for matching value")
	}
}

func TestRegexFilter_Match_NoMatch(t *testing.T) {
	f, _ := NewRegexFilter("msg", `^fatal`)
	line := makeRegexLogLine(map[string]string{"msg": "info: all good"})
	if f.Match(line) {
		t.Error("expected false for non-matching value")
	}
}
