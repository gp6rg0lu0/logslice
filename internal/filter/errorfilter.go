package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// ErrorFilter matches log lines where a specified field contains an error-like value.
// It checks for common error indicators: non-empty strings, known error keywords,
// or any non-null value depending on the mode.
type ErrorFilter struct {
	field    string
	keywords []string
}

// NewErrorFilter creates a filter that matches lines where the given field
// contains one of the provided keywords (case-insensitive). If no keywords
// are provided, any non-empty value in the field is considered a match.
func NewErrorFilter(field string, keywords []string) (*ErrorFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("errorfilter: field must not be empty")
	}
	lower := make([]string, len(keywords))
	for i, k := range keywords {
		if strings.TrimSpace(k) == "" {
			return nil, fmt.Errorf("errorfilter: keyword at index %d is blank", i)
		}
		lower[i] = strings.ToLower(k)
	}
	return &ErrorFilter{field: field, keywords: lower}, nil
}

// Field returns the field name this filter operates on.
func (f *ErrorFilter) Field() string { return f.field }

// Keywords returns the list of keywords used for matching.
func (f *ErrorFilter) Keywords() []string { return f.keywords }

// Match returns true if the log line's field value matches the filter criteria.
func (f *ErrorFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok || val == "" {
		return false
	}
	if len(f.keywords) == 0 {
		return true
	}
	lower := strings.ToLower(val)
	for _, kw := range f.keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}
