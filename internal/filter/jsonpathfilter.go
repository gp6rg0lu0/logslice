package filter

import (
	"fmt"
	"strings"

	"github.com/nicholasgasior/logslice/internal/parser"
)

// JSONPathFilter matches log lines where a dot-separated field path equals a target value.
type JSONPathFilter struct {
	path   string
	value  string
	parts  []string
}

// NewJSONPathFilter creates a filter that resolves a dot-separated path in nested JSON
// and checks whether the resolved value equals the given target.
func NewJSONPathFilter(path, value string) (*JSONPathFilter, error) {
	if path == "" {
		return nil, fmt.Errorf("jsonpath filter: path must not be empty")
	}
	if value == "" {
		return nil, fmt.Errorf("jsonpath filter: value must not be empty")
	}
	parts := strings.Split(path, ".")
	for _, p := range parts {
		if p == "" {
			return nil, fmt.Errorf("jsonpath filter: path segment must not be empty in %q", path)
		}
	}
	return &JSONPathFilter{path: path, value: value, parts: parts}, nil
}

// Path returns the dot-separated field path.
func (f *JSONPathFilter) Path() string { return f.path }

// Value returns the target value.
func (f *JSONPathFilter) Value() string { return f.value }

// Match returns true if the resolved nested value equals the target.
func (f *JSONPathFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	current := line.Fields()
	for i, part := range f.parts {
		val, ok := current[part]
		if !ok {
			return false
		}
		if i == len(f.parts)-1 {
			s, ok := val.(string)
			if !ok {
				return fmt.Sprintf("%v", val) == f.value
			}
			return s == f.value
		}
		nested, ok := val.(map[string]interface{})
		if !ok {
			return false
		}
		current = nested
	}
	return false
}
