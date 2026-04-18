package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// NotContainsFilter matches log lines where a field does NOT contain a substring.
type NotContainsFilter struct {
	field     string
	substring string
	caseSensitive bool
}

// NewNotContainsFilter creates a filter that rejects lines where field contains substring.
func NewNotContainsFilter(field, substring string, caseSensitive bool) (*NotContainsFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("notcontains: field must not be empty")
	}
	if substring == "" {
		return nil, fmt.Errorf("notcontains: substring must not be empty")
	}
	return &NotContainsFilter{field: field, substring: substring, caseSensitive: caseSensitive}, nil
}

// Match returns true if the field does not contain the substring.
func (f *NotContainsFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return true
	}
	if f.caseSensitive {
		return !strings.Contains(val, f.substring)
	}
	return !strings.Contains(strings.ToLower(val), strings.ToLower(f.substring))
}

// Field returns the field name.
func (f *NotContainsFilter) Field() string { return f.field }

// Substring returns the substring.
func (f *NotContainsFilter) Substring() string { return f.substring }
