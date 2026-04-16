package filter

import (
	"fmt"
	"sync"

	"github.com/nicholasgasior/logslice/internal/parser"
)

// TailFilter keeps a sliding window of the last N matched lines.
// It buffers lines and only emits them once the stream ends via Flush.
type TailFilter struct {
	mu  sync.Mutex
	buf []*parser.LogLine
	max int
}

// NewTailFilter creates a TailFilter that retains the last n lines.
func NewTailFilter(n int) (*TailFilter, error) {
	if n <= 0 {
		return nil, fmt.Errorf("tailfilter: n must be greater than zero, got %d", n)
	}
	return &TailFilter{max: n, buf: make([]*parser.LogLine, 0, n)}, nil
}

// Match always returns true but stores the line in the ring buffer.
func (f *TailFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.buf) >= f.max {
		f.buf = append(f.buf[1:], line)
	} else {
		f.buf = append(f.buf, line)
	}
	return false // suppress normal output; use Lines() to retrieve
}

// Lines returns the buffered tail lines and resets the buffer.
func (f *TailFilter) Lines() []*parser.LogLine {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]*parser.LogLine, len(f.buf))
	copy(out, f.buf)
	f.buf = f.buf[:0]
	return out
}

// Max returns the configured window size.
func (f *TailFilter) Max() int { return f.max }
