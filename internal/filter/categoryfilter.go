package filter

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// CategoryFilter matches log lines where a field's value belongs to one of the
// specified categories. Matching is case-insensitive.
type CategoryFilter struct {
	field      string
	categories map[string]struct{}
	raw        []string
}

// NewCategoryFilter creates a CategoryFilter that matches when the given field
// contains one of the provided category values (case-insensitive).
func NewCategoryFilter(field string, categories []string) (*CategoryFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("categoryfilter: field must not be empty")
	}
	if len(categories) == 0 {
		return nil, fmt.Errorf("categoryfilter: at least one category is required")
	}
	set := make(map[string]struct{}, len(categories))
	raw := make([]string, 0, len(categories))
	for _, c := range categories {
		trimmed := strings.TrimSpace(c)
		if trimmed == "" {
			return nil, fmt.Errorf("categoryfilter: category value must not be blank")
		}
		key := strings.ToLower(trimmed)
		if _, exists := set[key]; !exists {
			set[key] = struct{}{}
			raw = append(raw, trimmed)
		}
	}
	return &CategoryFilter{field: field, categories: set, raw: raw}, nil
}

// Field returns the field name this filter inspects.
func (f *CategoryFilter) Field() string { return f.field }

// Categories returns the accepted category values.
func (f *CategoryFilter) Categories() []string { return f.raw }

// Match returns true when the log line's field value (case-insensitive) is one
// of the configured categories.
func (f *CategoryFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	_, found := f.categories[strings.ToLower(val)]
	return found
}
