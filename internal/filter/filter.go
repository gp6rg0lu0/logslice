package filter

// Filter is the common interface implemented by all log filters.
type Filter interface {
	// Match returns true if the given log entry (represented as a map of
	// parsed fields) satisfies the filter condition.
	Match(entry map[string]string) bool
}

// Chain combines multiple filters with AND semantics: all filters must match.
type Chain struct {
	filters []Filter
}

// NewChain creates a Chain from the provided filters. Nil filters are ignored.
func NewChain(filters ...Filter) *Chain {
	var active []Filter
	for _, f := range filters {
		if f != nil {
			active = append(active, f)
		}
	}
	return &Chain{filters: active}
}

// Match returns true only when every filter in the chain matches the entry.
func (c *Chain) Match(entry map[string]string) bool {
	for _, f := range c.filters {
		if !f.Match(entry) {
			return false
		}
	}
	return true
}

// Len returns the number of active filters in the chain.
func (c *Chain) Len() int {
	return len(c.filters)
}
