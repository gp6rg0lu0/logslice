package filter

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// ExistsFilter matches log lines where a given field is present or absent.
type ExistsFilter struct {
	field     string
	mustExist bool
}

// NewExistsFilter returns a filter that passes lines where field exists
// when mustExist is true, or where field is absent when mustExist is false.
func NewExistsFilter(field string, mustExist bool) (*ExistsFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("existsfilter: field name must not be empty")
	}
	return &ExistsFilter{field: field, mustExist: mustExist}, nil
}

// Field returns the field name being checked.
func (f *ExistsFilter) Field() string { return f.field }

// MustExist returns whether the filter requires the field to be present.
func (f *ExistsFilter) MustExist() bool { return f.mustExist }

// Match returns true when the line's field presence matches the mustExist flag.
func (f *ExistsFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	_, ok := line.Get(f.field)
	return ok == f.mustExist
}
