package filter

// OrFilter matches a log line when at least one of its inner filters matches.
type OrFilter struct {
	filters []Filter
}

// NewOrFilter creates an OrFilter from the provided filters, ignoring nils.
// An OrFilter with no active filters never matches.
func NewOrFilter(filters ...Filter) *OrFilter {
	active := make([]Filter, 0, len(filters))
	for _, f := range filters {
		if f != nil {
			active = append(active, f)
		}
	}
	return &OrFilter{filters: active}
}

// Match returns true if at least one inner filter matches the line.
func (o *OrFilter) Match(line interface{ Get(string) (string, bool) }) bool {
	for _, f := range o.filters {
		if f.Match(line) {
			return true
		}
	}
	return false
}

// Len returns the number of active inner filters.
func (o *OrFilter) Len() int {
	return len(o.filters)
}
