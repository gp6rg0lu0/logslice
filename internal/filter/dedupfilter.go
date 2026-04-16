package filter

import (
	"fmt"
	"sync"

	"github.com/user/logslice/internal/parser"
)

// DedupFilter drops log lines whose field value has already been seen.
type DedupFilter struct {
	mu    sync.Mutex
	field string
	seen  map[string]struct{}
}

// NewDedupFilter returns a DedupFilter deduplicating on the given field.
func NewDedupFilter(field string) (*DedupFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("dedupfilter: field must not be empty")
	}
	return &DedupFilter{
		field: field,
		seen:  make(map[string]struct{}),
	}, nil
}

// Field returns the deduplication field name.
func (f *DedupFilter) Field() string { return f.field }

// Match returns true only the first time a given field value is encountered.
func (f *DedupFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, exists := f.seen[val]; exists {
		return false
	}
	f.seen[val] = struct{}{}
	return true
}
