package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeCountLogLine(msg string) *parser.LogLine {
	return parser.NewLogLine(map[string]interface{}{"msg": msg})
}

func TestNewCountFilter_InvalidMax(t *testing.T) {
	for _, max := range []int{0, -1, -100} {
		_, err := NewCountFilter(max)
		if err == nil {
			t.Errorf("expected error for max=%d", max)
		}
	}
}

func TestNewCountFilter_ValidMax(t *testing.T) {
	f, err := NewCountFilter(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Max() != 5 {
		t.Errorf("expected Max()=5, got %d", f.Max())
	}
}

func TestCountFilter_NilLine(t *testing.T) {
	f, _ := NewCountFilter(3)
	if f.Match(nil) {
		t.Error("expected nil line to not match")
	}
	if f.Seen() != 0 {
		t.Errorf("nil line should not increment seen, got %d", f.Seen())
	}
}

func TestCountFilter_MatchesUpToMax(t *testing.T) {
	f, _ := NewCountFilter(3)
	line := makeCountLogLine("hello")

	for i := 0; i < 3; i++ {
		if !f.Match(line) {
			t.Errorf("expected match on call %d", i+1)
		}
	}
	if f.Match(line) {
		t.Error("expected no match after max reached")
	}
}

func TestCountFilter_SeenIncrements(t *testing.T) {
	f, _ := NewCountFilter(10)
	line := makeCountLogLine("x")
	for i := 0; i < 4; i++ {
		f.Match(line)
	}
	if f.Seen() != 4 {
		t.Errorf("expected Seen()=4, got %d", f.Seen())
	}
}

func TestCountFilter_MaxOne(t *testing.T) {
	f, _ := NewCountFilter(1)
	line := makeCountLogLine("only one")
	if !f.Match(line) {
		t.Error("expected first match")
	}
	if f.Match(line) {
		t.Error("expected no match after first")
	}
}
