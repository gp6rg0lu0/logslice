package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeSamplerLogLine() *parser.LogLine {
	return parser.NewLogLine(map[string]string{"msg": "test"})
}

func TestNewSamplerFilter_InvalidN(t *testing.T) {
	for _, n := range []int{0, -1, -100} {
		_, err := NewSamplerFilter(n)
		if err == nil {
			t.Errorf("expected error for n=%d", n)
		}
	}
}

func TestNewSamplerFilter_ValidN(t *testing.T) {
	f, err := NewSamplerFilter(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.N() != 3 {
		t.Errorf("expected N=3, got %d", f.N())
	}
}

func TestSamplerFilter_NilLine(t *testing.T) {
	f, _ := NewSamplerFilter(2)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestSamplerFilter_EveryNth(t *testing.T) {
	f, _ := NewSamplerFilter(3)
	results := make([]bool, 9)
	for i := range results {
		results[i] = f.Match(makeSamplerLogLine())
	}
	expected := []bool{false, false, true, false, false, true, false, false, true}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("index %d: expected %v, got %v", i, expected[i], got)
		}
	}
}

func TestSamplerFilter_N1_PassesAll(t *testing.T) {
	f, _ := NewSamplerFilter(1)
	for i := 0; i < 5; i++ {
		if !f.Match(makeSamplerLogLine()) {
			t.Errorf("expected true at index %d", i)
		}
	}
}

func TestSamplerFilter_Reset(t *testing.T) {
	f, _ := NewSamplerFilter(2)
	f.Match(makeSamplerLogLine()) // 1 - false
	f.Match(makeSamplerLogLine()) // 2 - true
	f.Reset()
	if f.Match(makeSamplerLogLine()) { // 1 again - false
		t.Error("expected false after reset")
	}
	if !f.Match(makeSamplerLogLine()) { // 2 again - true
		t.Error("expected true after reset on second call")
	}
}
