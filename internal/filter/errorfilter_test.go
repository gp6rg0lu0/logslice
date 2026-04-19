package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeErrorLogLine(fields map[string]string) *parser.LogLine {
	m := make(map[string]interface{})
	for k, v := range fields {
		m[k] = v
	}
	return parser.NewLogLine(m)
}

func TestNewErrorFilter_EmptyField(t *testing.T) {
	_, err := NewErrorFilter("", nil)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewErrorFilter_BlankKeyword(t *testing.T) {
	_, err := NewErrorFilter("error", []string{"timeout", "  "})
	if err == nil {
		t.Fatal("expected error for blank keyword")
	}
}

func TestNewErrorFilter_Valid(t *testing.T) {
	f, err := NewErrorFilter("msg", []string{"timeout", "refused"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "msg" {
		t.Errorf("expected field 'msg', got %q", f.Field())
	}
	if len(f.Keywords()) != 2 {
		t.Errorf("expected 2 keywords, got %d", len(f.Keywords()))
	}
}

func TestErrorFilter_NilLine(t *testing.T) {
	f, _ := NewErrorFilter("error", nil)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestErrorFilter_MissingField(t *testing.T) {
	f, _ := NewErrorFilter("error", nil)
	line := makeErrorLogLine(map[string]string{"msg": "something"})
	if f.Match(line) {
		t.Error("expected false when field is missing")
	}
}

func TestErrorFilter_NoKeywords_AnyValue(t *testing.T) {
	f, _ := NewErrorFilter("error", nil)
	line := makeErrorLogLine(map[string]string{"error": "something went wrong"})
	if !f.Match(line) {
		t.Error("expected true for non-empty field with no keywords")
	}
}

func TestErrorFilter_KeywordMatch(t *testing.T) {
	f, _ := NewErrorFilter("msg", []string{"timeout", "refused"})
	cases := []struct {
		val  string
		want bool
	}{
		{"connection timeout exceeded", true},
		{"REFUSED by server", true},
		{"all good", false},
		{"", false},
	}
	for _, tc := range cases {
		line := makeErrorLogLine(map[string]string{"msg": tc.val})
		if got := f.Match(line); got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.val, got, tc.want)
		}
	}
}
