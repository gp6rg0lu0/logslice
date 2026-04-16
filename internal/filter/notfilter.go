package filter

// NotFilter negates another filter's Match result.
type NotFilter struct {
	inner Filter
}

// NewNotFilter wraps the given filter and inverts its result.
// Returns nil if inner is nil.
func NewNotFilter(inner Filter) *NotFilter {
	if inner == nil {
		return nil
	}
	return &NotFilter{inner: inner}
}

// Match returns true when the inner filter returns false, and vice versa.
func (n *NotFilter) Match(line interface{ Get(string) (string, bool) }) bool {
	return !n.inner.Match(line)
}

// Inner returns the wrapped filter.
func (n *NotFilter) Inner() Filter {
	return n.inner
}
