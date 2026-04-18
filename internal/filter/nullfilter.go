package filter

import "github.com/logslice/logslice/internal/parser"

// NullFilter matches log lines where a given field is either null/missing
// or explicitly set to a non-null value, depending on the mustBeNull flag.
type NullFilter struct {
	field     string
	mustBeNull bool
}

// NewNullFilter creates a filter that checks whether field is null (missing or
// empty string). If mustBeNull is true, the line matches when the field is
// absent or empty; if false, it matches when the field is present and non-empty.
func NewNullFilter(field string, mustBeNull bool) (*NullFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("nullfilter: field must not be empty")
	}
	return &NullFilter{field: field, mustBeNull: mustBeNull}, nil
}

// Field returns the field name being checked.
func (f *NullFilter) Field() string { return f.field }

// MustBeNull returns the null expectation flag.
func (f *NullFilter) MustBeNull() bool { return f.mustBeNull }

// Match returns true when the field's null state matches the filter expectation.
func (f *NullFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	isNull := !ok || val == ""
	return isNull == f.mustBeNull
}
