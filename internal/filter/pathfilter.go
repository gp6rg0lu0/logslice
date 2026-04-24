package filter

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// PathFilter matches log lines where a field value matches one of the given
// URL path prefixes or exact paths.
type PathFilter struct {
	field   string
	paths   []string
	exact   bool
}

// NewPathFilter creates a PathFilter that checks whether the value of field
// starts with (or exactly matches, when exact=true) any of the given paths.
// Returns an error if field is empty or no paths are provided.
func NewPathFilter(field string, paths []string, exact bool) (*PathFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("pathfilter: field must not be empty")
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("pathfilter: at least one path must be provided")
	}
	cleaned := make([]string, 0, len(paths))
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, fmt.Errorf("pathfilter: path must not be blank")
		}
		cleaned = append(cleaned, p)
	}
	return &PathFilter{field: field, paths: cleaned, exact: exact}, nil
}

// Field returns the field name this filter operates on.
func (f *PathFilter) Field() string { return f.field }

// Paths returns the configured path patterns.
func (f *PathFilter) Paths() []string { return f.paths }

// Exact returns whether the filter uses exact matching.
func (f *PathFilter) Exact() bool { return f.exact }

// Match returns true if the log line's field value matches any configured path.
func (f *PathFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	for _, p := range f.paths {
		if f.exact {
			if val == p {
				return true
			}
		} else {
			if strings.HasPrefix(val, p) {
				return true
			}
		}
	}
	return false
}
