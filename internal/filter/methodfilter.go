package filter

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// MethodFilter matches log lines where a named field contains one of the
// specified HTTP methods (GET, POST, PUT, DELETE, etc.).
type MethodFilter struct {
	field   string
	methods map[string]struct{}
}

// NewMethodFilter creates a MethodFilter that matches when the given field
// contains any of the provided HTTP methods. Methods are normalized to
// uppercase. At least one method must be provided.
func NewMethodFilter(field string, methods []string) (*MethodFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("methodfilter: field must not be empty")
	}
	if len(methods) == 0 {
		return nil, fmt.Errorf("methodfilter: at least one HTTP method must be specified")
	}

	set := make(map[string]struct{}, len(methods))
	for _, m := range methods {
		norm := strings.ToUpper(strings.TrimSpace(m))
		if norm == "" {
			return nil, fmt.Errorf("methodfilter: method must not be blank")
		}
		set[norm] = struct{}{}
	}

	return &MethodFilter{field: field, methods: set}, nil
}

// Match returns true when the log line's field value (uppercased) is one of
// the configured HTTP methods.
func (f *MethodFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	_, found := f.methods[strings.ToUpper(strings.TrimSpace(val))]
	return found
}

// Field returns the field name inspected by this filter.
func (f *MethodFilter) Field() string { return f.field }

// Methods returns the set of HTTP methods this filter accepts.
func (f *MethodFilter) Methods() []string {
	out := make([]string, 0, len(f.methods))
	for m := range f.methods {
		out = append(out, m)
	}
	return out
}
