package filter_test

import (
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makePrefixLogLine(fields map[string]string) *parser.LogLine {
	m := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		m[k] = v
	}
	return parser.NewLogLine(m)
}

func TestNewPrefixFilter_EmptyField(t *testing.T) {
	_, err := filter.NewPrefixFilter("", "ERR")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewPrefixFilter_EmptyPrefix(t *testing.T) {
	_, err := filter.NewPrefixFilter("msg", "")
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestPrefixFilter_Accessors(t *testing.T) {
	f, err := filter.NewPrefixFilter("msg", "ERR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "msg" {
		t.Errorf("expected field 'msg', got %q", f.Field())
	}
	if f.Prefix() != "ERR" {
		t.Errorf("expected prefix 'ERR', got %q", f.Prefix())
	}
}

func TestPrefixFilter_NilLine(t *testing.T) {
	f, _ := filter.NewPrefixFilter("msg", "ERR")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestPrefixFilter_MatchesPrefix(t *testing.T) {
	f, _ := filter.NewPrefixFilter("msg", "ERR")
	line := makePrefixLogLine(map[string]string{"msg": "ERR something failed"})
	if !f.Match(line) {
		t.Error("expected match for line with matching prefix")
	}
}

func TestPrefixFilter_NoMatch(t *testing.T) {
	f, _ := filter.NewPrefixFilter("msg", "ERR")
	line := makePrefixLogLine(map[string]string{"msg": "INFO all good"})
	if f.Match(line) {
		t.Error("expected no match for line without prefix")
	}
}

func TestPrefixFilter_MissingField(t *testing.T) {
	f, _ := filter.NewPrefixFilter("msg", "ERR")
	line := makePrefixLogLine(map[string]string{"level": "error"})
	if f.Match(line) {
		t.Error("expected no match when field is missing")
	}
}
