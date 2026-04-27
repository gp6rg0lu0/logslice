package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// ServiceFilter matches log lines where a named field equals one of the
// allowed service names. Matching is case-insensitive.
type ServiceFilter struct {
	field    string
	services map[string]struct{}
}

// NewServiceFilter creates a ServiceFilter that accepts log lines whose
// field value matches any of the provided service names.
//
// field must be non-empty. services must contain at least one non-blank entry.
func NewServiceFilter(field string, services []string) (*ServiceFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("servicefilter: field must not be empty")
	}
	if len(services) == 0 {
		return nil, fmt.Errorf("servicefilter: at least one service name is required")
	}

	set := make(map[string]struct{}, len(services))
	for _, s := range services {
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			return nil, fmt.Errorf("servicefilter: service name must not be blank")
		}
		set[strings.ToLower(trimmed)] = struct{}{}
	}

	return &ServiceFilter{field: field, services: set}, nil
}

// Field returns the log field inspected by this filter.
func (f *ServiceFilter) Field() string { return f.field }

// Services returns the set of allowed service names (lowercased).
func (f *ServiceFilter) Services() []string {
	out := make([]string, 0, len(f.services))
	for s := range f.services {
		out = append(out, s)
	}
	return out
}

// Match returns true when the log line's field value (case-insensitive) is
// one of the configured service names.
func (f *ServiceFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	_, found := f.services[strings.ToLower(val)]
	return found
}
