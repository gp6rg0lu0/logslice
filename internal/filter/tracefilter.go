package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// TraceFilter matches log lines whose trace ID field matches one of the given prefixes or exact values.
type TraceFilter struct {
	field  string
	values []string
	exact  bool
}

// NewTraceFilter creates a filter that matches when the given field equals (exact=true)
// or has a prefix matching one of the provided trace IDs.
func NewTraceFilter(field string, values []string, exact bool) (*TraceFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("tracefilter: field must not be empty")
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("tracefilter: at least one trace value required")
	}
	var cleaned []string
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			return nil, fmt.Errorf("tracefilter: trace values must not be blank")
		}
		cleaned = append(cleaned, v)
	}
	return &TraceFilter{field: field, values: cleaned, exact: exact}, nil
}

func (f *TraceFilter) Field() string    { return f.field }
func (f *TraceFilter) Values() []string { return f.values }
func (f *TraceFilter) Exact() bool      { return f.exact }

func (f *TraceFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(f.field)
	if !ok {
		return false
	}
	for _, tv := range f.values {
		if f.exact {
			if v == tv {
				return true
			}
		} else {
			if strings.HasPrefix(v, tv) {
				return true
			}
		}
	}
	return false
}
