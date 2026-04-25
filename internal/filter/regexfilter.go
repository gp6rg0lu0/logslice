package filter

import (
	"fmt"
	"regexp"

	"github.com/user/logslice/internal/parser"
)

// RegexFilter matches log lines where a field's value matches a regular expression.
type RegexFilter struct {
	field string
	pattern *regexp.Regexp
}

// NewRegexFilter creates a RegexFilter that checks whether the given field
// matches the provided regular expression pattern.
func NewRegexFilter(field, pattern string) (*RegexFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("regexfilter: field name must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("regexfilter: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("regexfilter: invalid pattern %q: %w", pattern, err)
	}
	return &RegexFilter{field: field, pattern: re}, nil
}

// Match returns true when the log line's field value matches the regex.
func (f *RegexFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	return f.pattern.MatchString(val)
}

// String returns a human-readable description of the filter, useful for
// logging and debugging purposes.
func (f *RegexFilter) String() string {
	return fmt.Sprintf("RegexFilter{field: %q, pattern: %q}", f.field, f.pattern.String())
}

// Field returns the field name this filter operates on.
func (f *RegexFilter) Field() string { return f.field }

// Pattern returns the compiled regular expression.
func (f *RegexFilter) Pattern() *regexp.Regexp { return f.pattern }
