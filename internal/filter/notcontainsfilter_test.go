package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeNotContainsLogLine(key, val string) *parser.LogLine {
	return parser.NewLogLine(map[string]string{key: val})
}

func TestNewNotContainsFilter_EmptyField(t *testing.T) {
	_, err := NewNotContainsFilter("", "foo", true)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewNotContainsFilter_EmptySubstring(t *testing.T) {
	_, err := NewNotContainsFilter("msg", "", true)
	if err == nil {
		t.Fatal("expected error for empty substring")
	}
}

func TestNotContainsFilter_Accessors(t *testing.T) {
	f, _ := NewNotContainsFilter("msg", "error", true)
	if f.Field() != "msg" {
		t.Errorf("expected msg, got %s", f.Field())
	}
	if f.Substring() != "error" {
		t.Errorf("expected error, got %s", f.Substring())
	}
}

func TestNotContainsFilter_NilLine(t *testing.T) {
	f, _ := NewNotContainsFilter("msg", "error", true)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestNotContainsFilter_FieldAbsent(t *testing.T) {
	f, _ := NewNotContainsFilter("msg", "error", true)
	line := makeNotContainsLogLine("other", "something")
	if !f.Match(line) {
		t.Error("expected true when field is absent")
	}
}

func TestNotContainsFilter_DoesNotContain(t *testing.T) {
	f, _ := NewNotContainsFilter("msg", "error", true)
	line := makeNotContainsLogLine("msg", "all good")
	if !f.Match(line) {
		t.Error("expected true when field does not contain substring")
	}
}

func TestNotContainsFilter_Contains(t *testing.T) {
	f, _ := NewNotContainsFilter("msg", "error", true)
	line := makeNotContainsLogLine("msg", "an error occurred")
	if f.Match(line) {
		t.Error("expected false when field contains substring")
	}
}

func TestNotContainsFilter_CaseInsensitive(t *testing.T) {
	f, _ := NewNotContainsFilter("msg", "error", false)
	line := makeNotContainsLogLine("msg", "An ERROR occurred")
	if f.Match(line) {
		t.Error("expected false for case-insensitive match")
	}
}
