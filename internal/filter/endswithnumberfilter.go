package filter

import (
	"fmt"
	"unicode"

	"github.com/user/logslice/internal/parser"
)

// EndsWithNumberFilter matches log lines where a field's value ends with a numeric digit.
type EndsWithNumberFilter struct {
	field  string
	invert bool
}

// NewEndsWithNumberFilter creates a filter that checks if the given field ends with a digit.
// If invert is true, it matches lines where the field does NOT end with a digit.
func NewEndsWithNumberFilter(field string, invert bool) (*EndsWithNumberFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("endswithnumber: field must not be empty")
	}
	return &EndsWithNumberFilter{field: field, invert: invert}, nil
}

// Match returns true if the field value ends with a numeric digit (inverted if configured).
func (f *EndsWithNumberFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok || val == "" {
		return false
	}
	runes := []rune(val)
	matched := unicode.IsDigit(runes[len(runes)-1])
	if f.invert {
		return !matched
	}
	return matched
}

// Field returns the field name being checked.
func (f *EndsWithNumberFilter) Field() string { return f.field }

// Invert returns whether the filter is inverted.
func (f *EndsWithNumberFilter) Invert() bool { return f.invert }
