package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// ContainsFilter matches log lines where a field contains a substring.
type ContainsFilter struct {
	field     string
	substring string
	caseFold  bool
}

// NewContainsFilter creates a filter that matches when field contains substring.
// If caseFold is true, comparison is case-insensitive.
func NewContainsFilter(field, substring string, caseFold bool) (*ContainsFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("containsfilter: field must not be empty")
	}
	if substring == "" {
		return nil, fmt.Errorf("containsfilter: substring must not be empty")
	}
	return &ContainsFilter{field: field, substring: substring, caseFold: caseFold}, nil
}

// Field returns the field name being checked.
func (f *ContainsFilter) Field() string { return f.field }

// Substring returns the substring being searched for.
func (f *ContainsFilter) Substring() string { return f.substring }

// CaseFold returns whether matching is case-insensitive.
func (f *ContainsFilter) CaseFold() bool { return f.caseFold }

// Match returns true if the field value contains the configured substring.
func (f *ContainsFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	if f.caseFold {
		return strings.Contains(strings.ToLower(val), strings.ToLower(f.substring))
	}
	return strings.Contains(val, f.substring)
}
