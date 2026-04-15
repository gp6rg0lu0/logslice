package filter

import (
	"fmt"
	"regexp"
)

// FieldFilter matches log entries where a specific field matches a pattern.
type FieldFilter struct {
	field   string
	pattern *regexp.Regexp
}

// NewFieldFilter creates a FieldFilter that matches entries where the given
// field value matches the provided regex pattern.
func NewFieldFilter(field, pattern string) (*FieldFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("field name must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern %q: %w", pattern, err)
	}
	return &FieldFilter{field: field, pattern: re}, nil
}

// Match returns true if the entry map contains the field and its string value
// matches the compiled pattern. Non-string field values are converted via
// fmt.Sprintf before matching.
func (f *FieldFilter) Match(entry map[string]interface{}) bool {
	val, ok := entry[f.field]
	if !ok {
		return false
	}
	var s string
	switch v := val.(type) {
	case string:
		s = v
	default:
		s = fmt.Sprintf("%v", v)
	}
	return f.pattern.MatchString(s)
}

// Field returns the field name this filter operates on.
func (f *FieldFilter) Field() string {
	return f.field
}

// Pattern returns the compiled regex pattern as a string.
func (f *FieldFilter) Pattern() string {
	return f.pattern.String()
}
