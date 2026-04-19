package filter

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/logslice/internal/parser"
)

// HostnameFilter matches log lines where a field value matches one of the
// provided hostname glob patterns (supports * as wildcard suffix/prefix).
type HostnameFilter struct {
	field    string
	patterns []string
}

// NewHostnameFilter creates a HostnameFilter that matches when the given field
// value matches any of the provided hostname patterns.
func NewHostnameFilter(field string, patterns []string) (*HostnameFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("hostname filter: field must not be empty")
	}
	if len(patterns) == 0 {
		return nil, fmt.Errorf("hostname filter: at least one pattern is required")
	}
	for _, p := range patterns {
		if p == "" {
			return nil, fmt.Errorf("hostname filter: pattern must not be empty")
		}
	}
	return &HostnameFilter{field: field, patterns: patterns}, nil
}

// Field returns the field name being filtered.
func (f *HostnameFilter) Field() string { return f.field }

// Patterns returns the hostname patterns.
func (f *HostnameFilter) Patterns() []string { return f.patterns }

// Match returns true if the log line's field value matches any pattern.
func (f *HostnameFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	v := strings.ToLower(val)
	for _, p := range f.patterns {
		pat := strings.ToLower(p)
		if matchHostnamePattern(pat, v) {
			return true
		}
	}
	return false
}

func matchHostnamePattern(pattern, value string) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[1:]
		return strings.HasSuffix(value, suffix) || value == pattern[2:]
	}
	if strings.HasSuffix(pattern, ".*") {
		prefix := pattern[:len(pattern)-2]
		return value == prefix || strings.HasPrefix(value, prefix+".")
	}
	return value == pattern
}
