package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// TimestampFilter matches log lines where a named field parses as a time
// and falls within [min, max]. Either bound may be zero to indicate unbounded.
type TimestampFilter struct {
	field  string
	min    time.Time
	max    time.Time
	layout string
}

// NewTimestampFilter creates a filter that checks a field's parsed timestamp.
// layout is the time.Parse layout string. min/max are optional (zero = unbounded).
func NewTimestampFilter(field, layout string, min, max time.Time) (*TimestampFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("timestampfilter: field must not be empty")
	}
	if layout == "" {
		return nil, fmt.Errorf("timestampfilter: layout must not be empty")
	}
	if !min.IsZero() && !max.IsZero() && max.Before(min) {
		return nil, fmt.Errorf("timestampfilter: max must not be before min")
	}
	return &TimestampFilter{field: field, layout: layout, min: min, max: max}, nil
}

func (f *TimestampFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(f.field)
	if !ok {
		return false
	}
	t, err := time.Parse(f.layout, v)
	if err != nil {
		return false
	}
	if !f.min.IsZero() && t.Before(f.min) {
		return false
	}
	if !f.max.IsZero() && t.After(f.max) {
		return false
	}
	return true
}

func (f *TimestampFilter) Field() string  { return f.field }
func (f *TimestampFilter) Layout() string { return f.layout }
func (f *TimestampFilter) Min() time.Time { return f.min }
func (f *TimestampFilter) Max() time.Time { return f.max }
