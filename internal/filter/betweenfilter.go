package filter

import (
	"fmt"
	"strconv"
)

// BetweenFilter matches log lines where a numeric field value falls within [Min, Max].
type BetweenFilter struct {
	field string
	min   float64
	max   float64
}

// NewBetweenFilter creates a BetweenFilter for the given field and inclusive range.
func NewBetweenFilter(field string, min, max float64) (*BetweenFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("betweenfilter: field must not be empty")
	}
	if min > max {
		return nil, fmt.Errorf("betweenfilter: min %.4g must not exceed max %.4g", min, max)
	}
	return &BetweenFilter{field: field, min: min, max: max}, nil
}

// Field returns the field name being filtered.
func (f *BetweenFilter) Field() string { return f.field }

// Min returns the lower bound.
func (f *BetweenFilter) Min() float64 { return f.min }

// Max returns the upper bound.
func (f *BetweenFilter) Max() float64 { return f.max }

// Match returns true if the field value is a number within [Min, Max].
func (f *BetweenFilter) Match(line interface{ Get(string) (string, bool) }) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(f.field)
	if !ok {
		return false
	}
	n, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return false
	}
	return n >= f.min && n <= f.max
}
