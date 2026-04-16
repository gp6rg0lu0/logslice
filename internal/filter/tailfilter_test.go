package filter

import (
	"testing"

	"github.com/nicholasgasior/logslice/internal/parser"
)

func makeTailLogLine(msg string) *parser.LogLine {
	return parser.NewLogLine(map[string]interface{}{"msg": msg})
}

func TestNewTailFilter_InvalidN(t *testing.T) {
	_, err := NewTailFilter(0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
	_, err = NewTailFilter(-5)
	if err == nil {
		t.Fatal("expected error for n=-5")
	}
}

func TestNewTailFilter_ValidN(t *testing.T) {
	f, err := NewTailFilter(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Max() != 3 {
		t.Errorf("expected max=3, got %d", f.Max())
	}
}

func TestTailFilter_NilLine(t *testing.T) {
	f, _ := NewTailFilter(3)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestTailFilter_BuffersLastN(t *testing.T) {
	f, _ := NewTailFilter(3)
	for i := 0; i < 5; i++ {
		f.Match(makeTailLogLine(string(rune('a' + i))))
	}
	lines := f.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	expected := []string{"c", "d", "e"}
	for i, l := range lines {
		if v, _ := l.Get("msg"); v != expected[i] {
			t.Errorf("line %d: expected %q got %q", i, expected[i], v)
		}
	}
}

func TestTailFilter_LinesResetsBuffer(t *testing.T) {
	f, _ := NewTailFilter(3)
	f.Match(makeTailLogLine("x"))
	f.Lines()
	if got := f.Lines(); len(got) != 0 {
		t.Errorf("expected empty buffer after Lines(), got %d", len(got))
	}
}

func TestTailFilter_MatchReturnsFalse(t *testing.T) {
	f, _ := NewTailFilter(2)
	line := makeTailLogLine("hello")
	if f.Match(line) {
		t.Error("Match should always return false")
	}
}
