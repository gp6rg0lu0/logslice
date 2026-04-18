package filter

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// LenFilter matches log lines where the length of a field's value
// satisfies a min/max constraint.
type LenFilter struct {
	field  string
	minLen int
	maxLen int
}

// NewLenFilter returns a LenFilter that matches when len(field) is in [minLen, maxLen].
// Pass -1 for maxLen to indicate no upper bound.
func NewLenFilter(field string, minLen, maxLen int) (*LenFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("lenfilter: field must not be empty")
	}
	if minLen < 0 {
		return nil, fmt.Errorf("lenfilter: minLen must be >= 0")
	}
	if maxLen != -1 && maxLen < minLen {
		return nil, fmt.Errorf("lenfilter: maxLen must be >= minLen or -1")
	}
	return &LenFilter{field: field, minLen: minLen, maxLen: maxLen}, nil
}

// Match returns true when the field value length is within [minLen, maxLen].
func (f *LenFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(f.field)
	if !ok {
		return false
	}
	l := len(v)
	if l < f.minLen {
		return false
	}
	if f.maxLen != -1 && l > f.maxLen {
		return false
	}
	return true
}

// Field returns the field name being checked.
func (f *LenFilter) Field() string { return f.field }

// MinLen returns the minimum length.
func (f *LenFilter) MinLen() int { return f.minLen }

// MaxLen returns the maximum length (-1 means unbounded).
func (f *LenFilter) MaxLen() int { return f.maxLen }
