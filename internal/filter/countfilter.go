package filter

import (
	"fmt"
	"sync/atomic"

	"github.com/user/logslice/internal/parser"
)

// CountFilter matches only the first N log lines that pass through it.
type CountFilter struct {
	max   int64
	count atomic.Int64
}

// NewCountFilter returns a Filter that matches at most max lines.
// Returns an error if max is less than 1.
func NewCountFilter(max int) (*CountFilter, error) {
	if max < 1 {
		return nil, fmt.Errorf("countfilter: max must be >= 1, got %d", max)
	}
	return &CountFilter{max: int64(max)}, nil
}

// Match returns true until max lines have been matched, then always false.
func (f *CountFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	next := f.count.Add(1)
	return next <= f.max
}

// Max returns the configured maximum count.
func (f *CountFilter) Max() int64 {
	return f.max
}

// Seen returns how many lines have been evaluated so far.
func (f *CountFilter) Seen() int64 {
	return f.count.Load()
}
