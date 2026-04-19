package filter

import (
	"crypto/md5"
	"fmt"
	"sync"

	"github.com/yourorg/logslice/internal/parser"
)

// HashFilter deduplicates log lines by computing a hash over one or more fields.
// Unlike DedupFilter (which tracks a single field value), HashFilter hashes a
// composite key built from multiple fields, making it suitable for multi-field
// deduplication.
type HashFilter struct {
	fields []string
	mu     sync.Mutex
	seen   map[string]struct{}
}

// NewHashFilter returns a HashFilter that deduplicates based on the combined
// values of the given fields. At least one field must be provided.
func NewHashFilter(fields []string) (*HashFilter, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("hashfilter: at least one field is required")
	}
	for _, f := range fields {
		if f == "" {
			return nil, fmt.Errorf("hashfilter: field name must not be empty")
		}
	}
	return &HashFilter{
		fields: fields,
		seen:   make(map[string]struct{}),
	}, nil
}

// Fields returns the fields used to compute the hash.
func (h *HashFilter) Fields() []string { return h.fields }

// Match returns true the first time a given field-value combination is seen,
// and false for all subsequent occurrences.
func (h *HashFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}

	hasher := md5.New()
	for _, f := range h.fields {
		v, _ := line.Get(f)
		fmt.Fprintf(hasher, "%s=%s;", f, v)
	}
	key := fmt.Sprintf("%x", hasher.Sum(nil))

	h.mu.Lock()
	defer h.mu.Unlock()
	if _, exists := h.seen[key]; exists {
		return false
	}
	h.seen[key] = struct{}{}
	return true
}
