package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// RequestIDFilter matches log lines where a field value matches one or more
// request ID prefixes or exact values.
type RequestIDFilter struct {
	field  string
	values []string
	exact  bool
}

// NewRequestIDFilter creates a filter that matches log lines where the given
// field contains one of the specified request ID values. If exact is true,
// the full value must match; otherwise prefix matching is used.
func NewRequestIDFilter(field string, values []string, exact bool) (*RequestIDFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("requestid filter: field must not be empty")
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("requestid filter: at least one request ID value is required")
	}
	var cleaned []string
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			return nil, fmt.Errorf("requestid filter: request ID values must not be blank")
		}
		cleaned = append(cleaned, v)
	}
	return &RequestIDFilter{field: field, values: cleaned, exact: exact}, nil
}

// Field returns the log field inspected by this filter.
func (f *RequestIDFilter) Field() string { return f.field }

// Values returns the request ID values used for matching.
func (f *RequestIDFilter) Values() []string { return f.values }

// Exact returns true when exact matching is required.
func (f *RequestIDFilter) Exact() bool { return f.exact }

// Match returns true if the log line's field value matches any of the
// configured request ID values.
func (f *RequestIDFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(f.field)
	if !ok {
		return false
	}
	for _, id := range f.values {
		if f.exact {
			if v == id {
				return true
			}
		} else {
			if strings.HasPrefix(v, id) {
				return true
			}
		}
	}
	return false
}
