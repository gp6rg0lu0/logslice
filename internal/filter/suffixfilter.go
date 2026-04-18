package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// SuffixFilter matches log lines where a field value ends with a given suffix.
type SuffixFilter struct {
	field  string
	suffix string
}

// NewSuffixFilter creates a SuffixFilter for the given field and suffix.
func NewSuffixFilter(field, suffix string) (*SuffixFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("suffixfilter: field must not be empty")
	}
	if suffix == "" {
		return nil, fmt.Errorf("suffixfilter: suffix must not be empty")
	}
	return &SuffixFilter{field: field, suffix: suffix}, nil
}

// Field returns the field name used by this filter.
func (f *SuffixFilter) Field() string { return f.field }

// Suffix returns the suffix used by this filter.
func (f *SuffixFilter) Suffix() string { return f.suffix }

// Match returns true if the field value ends with the configured suffix.
func (f *SuffixFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	return strings.HasSuffix(val, f.suffix)
}
