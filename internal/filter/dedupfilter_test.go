package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeDedupLogLine(fields map[string]string) *parser.LogLine {
	raw := make(map[string]interface{})
	for k, v := range fields {
		raw[k] = v
	}
	return parser.NewLogLine(raw)
}

func TestNewDedupFilter_EmptyField(t *testing.T) {
	_, err := NewDedupFilter("")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewDedupFilter_ValidField(t *testing.T) {
	f, err := NewDedupFilter("request_id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "request_id" {
		t.Errorf("expected field request_id, got %s", f.Field())
	}
}

func TestDedupFilter_NilLine(t *testing.T) {
	f, _ := NewDedupFilter("id")
	if f.Match(nil) {
		t.Error("expected nil line to not match")
	}
}

func TestDedupFilter_MissingField(t *testing.T) {
	f, _ := NewDedupFilter("id")
	line := makeDedupLogLine(map[string]string{"msg": "hello"})
	if f.Match(line) {
		t.Error("expected line missing field to not match")
	}
}

func TestDedupFilter_UniqueLines(t *testing.T) {
	f, _ := NewDedupFilter("request_id")
	l1 := makeDedupLogLine(map[string]string{"request_id": "abc"})
	l2 := makeDedupLogLine(map[string]string{"request_id": "def"})
	if !f.Match(l1) {
		t.Error("expected first unique line to match")
	}
	if !f.Match(l2) {
		t.Error("expected second unique line to match")
	}
}

func TestDedupFilter_DuplicateLine(t *testing.T) {
	f, _ := NewDedupFilter("request_id")
	l1 := makeDedupLogLine(map[string]string{"request_id": "abc"})
	l2 := makeDedupLogLine(map[string]string{"request_id": "abc"})
	if !f.Match(l1) {
		t.Error("expected first line to match")
	}
	if f.Match(l2) {
		t.Error("expected duplicate line to not match")
	}
}
