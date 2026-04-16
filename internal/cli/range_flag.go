package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// RangeValue holds a parsed --range flag value in the form "field:min:max".
type RangeValue struct {
	Field string
	Min   float64
	Max   float64
}

// ParseRangeFlag parses a string of the form "field:min:max" into a RangeValue.
func ParseRangeFlag(s string) (RangeValue, error) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) != 3 {
		return RangeValue{}, fmt.Errorf("range flag %q must be in format field:min:max", s)
	}
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return RangeValue{}, fmt.Errorf("range flag: field name must not be empty")
	}
	min, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return RangeValue{}, fmt.Errorf("range flag: invalid min %q: %w", parts[1], err)
	}
	max, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err != nil {
		return RangeValue{}, fmt.Errorf("range flag: invalid max %q: %w", parts[2], err)
	}
	if min > max {
		return RangeValue{}, fmt.Errorf("range flag: min (%v) must be <= max (%v)", min, max)
	}
	return RangeValue{Field: field, Min: min, Max: max}, nil
}
