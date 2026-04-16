package filter

import "fmt"

// RangeFilter matches log lines where a numeric field falls within [Min, Max].
type RangeFilter struct {
	field string
	min   float64
	max   float64
}

// NewRangeFilter creates a RangeFilter for the given field and inclusive range.
func NewRangeFilter(field string, min, max float64) (*RangeFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("rangefilter: field must not be empty")
	}
	if min > max {
		return nil, fmt.Errorf("rangefilter: min (%v) must be <= max (%v)", min, max)
	}
	return &RangeFilter{field: field, min: min, max: max}, nil
}

// Field returns the field name being filtered.
func (r *RangeFilter) Field() string { return r.field }

// Min returns the lower bound of the range.
func (r *RangeFilter) Min() float64 { return r.min }

// Max returns the upper bound of the range.
func (r *RangeFilter) Max() float64 { return r.max }

// Match returns true if the log line's field value is a number within [min, max].
func (r *RangeFilter) Match(line interface{ Get(string) (string, bool) }) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(r.field)
	if !ok {
		return false
	}
	var num float64
	if _, err := fmt.Sscanf(v, "%g", &num); err != nil {
		return false
	}
	return num >= r.min && num <= r.max
}
