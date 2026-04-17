package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// PrefixFilter matches log lines where a field value starts with a given prefix.
type PrefixFilter struct {
	field  string
	prefix string
}

// NewPrefixFilter creates a PrefixFilter for the given field and prefix.
// Returns an error if field or prefix is empty.
func NewPrefixFilter(field, prefix string) (*PrefixFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("prefixfilter: field must not be empty")
	}
	if prefix == "" {
		return nil, fmt.Errorf("prefixfilter: prefix must not be empty")
	}
	return &PrefixFilter{field: field, prefix: prefix}, nil
}

// Match returns true if the line's field value starts with the configured prefix.
func (f *PrefixFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	return strings.HasPrefix(val, f.prefix)
}

// Field returns the field name used by this filter.
func (f *PrefixFilter) Field() string { return f.field }

// Prefix returns the prefix string used by this filter.
func (f *PrefixFilter) Prefix() string { return f.prefix }
