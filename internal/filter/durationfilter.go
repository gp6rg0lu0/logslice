package filter

import (
	"fmt"
	"time"

	"github.com/user/logslice/internal/parser"
)

// DurationFilter matches log lines where a duration field falls within [min, max].
type DurationFilter struct {
	field string
	min   time.Duration
	max   time.Duration
}

// NewDurationFilter creates a filter that matches lines where the named field,
// parsed as a Go duration string, is between min and max (inclusive).
func NewDurationFilter(field, minStr, maxStr string) (*DurationFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("durationfilter: field name must not be empty")
	}
	min, err := time.ParseDuration(minStr)
	if err != nil {
		return nil, fmt.Errorf("durationfilter: invalid min duration %q: %w", minStr, err)
	}
	max, err := time.ParseDuration(maxStr)
	if err != nil {
		return nil, fmt.Errorf("durationfilter: invalid max duration %q: %w", maxStr, err)
	}
	if max < min {
		return nil, fmt.Errorf("durationfilter: max (%s) must be >= min (%s)", maxStr, minStr)
	}
	return &DurationFilter{field: field, min: min, max: max}, nil
}

func (f *DurationFilter) Field() string        { return f.field }
func (f *DurationFilter) Min() time.Duration   { return f.min }
func (f *DurationFilter) Max() time.Duration   { return f.max }

// Match returns true if the field value parses as a duration within [min, max].
func (f *DurationFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return false
	}
	return d >= f.min && d <= f.max
}
