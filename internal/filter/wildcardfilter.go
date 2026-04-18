package filter

import (
	"fmt"
	"path"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// WildcardFilter matches log lines where a field value matches a glob pattern.
type WildcardFilter struct {
	field   string
	pattern string
}

// NewWildcardFilter creates a WildcardFilter for the given field and glob pattern.
// Returns an error if field or pattern is empty, or if the pattern is invalid.
func NewWildcardFilter(field, pattern string) (*WildcardFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("wildcard filter: field must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("wildcard filter: pattern must not be empty")
	}
	// Validate pattern by doing a test match.
	_, err := path.Match(pattern, "")
	if err != nil {
		return nil, fmt.Errorf("wildcard filter: invalid pattern %q: %w", pattern, err)
	}
	return &WildcardFilter{field: field, pattern: pattern}, nil
}

// Field returns the field name being matched.
func (f *WildcardFilter) Field() string { return f.field }

// Pattern returns the glob pattern.
func (f *WildcardFilter) Pattern() string { return f.pattern }

// Match returns true if the line's field value matches the glob pattern.
func (f *WildcardFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	matched, err := path.Match(strings.ToLower(f.pattern), strings.ToLower(val))
	if err != nil {
		return false
	}
	return matched
}
