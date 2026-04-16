package filter

import (
	"fmt"
	"sync/atomic"

	"github.com/user/logslice/internal/parser"
)

// SamplerFilter passes every Nth log line.
type SamplerFilter struct {
	n       int
	counter atomic.Int64
}

// NewSamplerFilter creates a filter that passes every nth line.
// n must be >= 1.
func NewSamplerFilter(n int) (*SamplerFilter, error) {
	if n < 1 {
		return nil, fmt.Errorf("sampler: n must be >= 1, got %d", n)
	}
	return &SamplerFilter{n: n}, nil
}

// Match returns true for every nth line.
func (f *SamplerFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	c := f.counter.Add(1)
	return c%int64(f.n) == 0
}

// N returns the sampling interval.
func (f *SamplerFilter) N() int {
	return f.n
}

// Reset resets the internal counter.
func (f *SamplerFilter) Reset() {
	f.counter.Store(0)
}
