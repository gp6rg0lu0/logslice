package filter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// LatencyFilter matches log lines where a numeric field (in milliseconds)
// falls within an optional [min, max] range.
type LatencyFilter struct {
	field string
	min   float64
	max   float64
	hasMin bool
	hasMax bool
}

// NewLatencyFilter creates a LatencyFilter. rangeStr is "min:max", "min:", or ":max".
func NewLatencyFilter(field, rangeStr string) (*LatencyFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("latency filter: field must not be empty")
	}
	parts := strings.SplitN(rangeStr, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("latency filter: range must be in format min:max")
	}
	f := &LatencyFilter{field: field}
	if parts[0] != "" {
		v, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return nil, fmt.Errorf("latency filter: invalid min %q: %w", parts[0], err)
		}
		f.min = v
		f.hasMin = true
	}
	if parts[1] != "" {
		v, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("latency filter: invalid max %q: %w", parts[1], err)
		}
		f.max = v
		f.hasMax = true
	}
	if f.hasMin && f.hasMax && f.min > f.max {
		return nil, fmt.Errorf("latency filter: min %.2f exceeds max %.2f", f.min, f.max)
	}
	return f, nil
}

func (f *LatencyFilter) Field() string  { return f.field }
func (f *LatencyFilter) Min() float64   { return f.min }
func (f *LatencyFilter) Max() float64   { return f.max }

func (f *LatencyFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return false
	}
	if f.hasMin && v < f.min {
		return false
	}
	if f.hasMax && v > f.max {
		return false
	}
	return true
}
