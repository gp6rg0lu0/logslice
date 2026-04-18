package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// MultiValueFilter matches log lines where a field's value is one of a set of allowed values.
type MultiValueFilter struct {
	field  string
	values map[string]struct{}
}

// NewMultiValueFilter creates a filter that matches when field equals any of the given values.
// At least one value must be provided.
func NewMultiValueFilter(field string, values []string) (*MultiValueFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("multivalue filter: field must not be empty")
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("multivalue filter: at least one value required")
	}
	set := make(map[string]struct{}, len(values))
	for _, v := range values {
		set[strings.TrimSpace(v)] = struct{}{}
	}
	return &MultiValueFilter{field: field, values: set}, nil
}

// Match returns true if the line's field value is in the allowed set.
func (f *MultiValueFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(f.field)
	if !ok {
		return false
	}
	_, found := f.values[v]
	return found
}

// Field returns the field name being filtered.
func (f *MultiValueFilter) Field() string { return f.field }

// Values returns the allowed values as a sorted slice.
func (f *MultiValueFilter) Values() []string {
	out := make([]string, 0, len(f.values))
	for v := range f.values {
		out = append(out, v)
	}
	return out
}
