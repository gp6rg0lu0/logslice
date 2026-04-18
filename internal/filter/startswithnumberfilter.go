package filter

import (
	"fmt"
	"unicode"

	"github.com/user/logslice/internal/parser"
)

// StartsWithNumberFilter matches log lines where a field value starts with a digit.
type StartsWithNumberFilter struct {
	field  string
	invert bool
}

// NewStartsWithNumberFilter returns a filter that matches lines where the given
// field's value starts with a numeric digit. If invert is true, lines that do
// NOT start with a digit are matched instead.
func NewStartsWithNumberFilter(field string, invert bool) (*StartsWithNumberFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("startswithnumber: field must not be empty")
	}
	return &StartsWithNumberFilter{field: field, invert: invert}, nil
}

// Field returns the field name this filter inspects.
func (f *StartsWithNumberFilter) Field() string { return f.field }

// Invert returns whether the match logic is inverted.
func (f *StartsWithNumberFilter) Invert() bool { return f.invert }

// Match returns true when the field value starts with a digit (or the inverse).
func (f *StartsWithNumberFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(f.field)
	if !ok || v == "" {
		return f.invert
	}
	starts := unicode.IsDigit(rune(v[0]))
	if f.invert {
		return !starts
	}
	return starts
}
