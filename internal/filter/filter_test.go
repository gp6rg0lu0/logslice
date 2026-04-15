package filter

import (
	"testing"
)

// alwaysFilter is a test helper that always returns the given value.
type alwaysFilter struct{ result bool }

func (a alwaysFilter) Match(_ map[string]string) bool { return a.result }

func TestNewChain_IgnoresNilFilters(t *testing.T) {
	c := NewChain(nil, alwaysFilter{true}, nil)
	if c.Len() != 1 {
		t.Errorf("expected 1 active filter, got %d", c.Len())
	}
}

func TestChain_AllMatch(t *testing.T) {
	c := NewChain(alwaysFilter{true}, alwaysFilter{true})
	if !c.Match(map[string]string{}) {
		t.Error("expected chain to match when all filters match")
	}
}

func TestChain_OneFails(t *testing.T) {
	c := NewChain(alwaysFilter{true}, alwaysFilter{false})
	if c.Match(map[string]string{}) {
		t.Error("expected chain to fail when one filter does not match")
	}
}

func TestChain_EmptyMatchesAll(t *testing.T) {
	c := NewChain()
	if !c.Match(map[string]string{"foo": "bar"}) {
		t.Error("empty chain should match every entry")
	}
}

func TestChain_Len(t *testing.T) {
	c := NewChain(alwaysFilter{true}, alwaysFilter{false}, alwaysFilter{true})
	if c.Len() != 3 {
		t.Errorf("expected Len 3, got %d", c.Len())
	}
}
