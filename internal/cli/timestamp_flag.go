package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/logslice/internal/filter"
)

// ParseTimestampFlag parses a flag value of the form "field:layout:min,max"
// where min and/or max may be omitted (empty string = unbounded).
// Example: "created_at:2006-01-02T15:04:05:2024-01-01T00:00:00,2024-12-31T23:59:59"
func ParseTimestampFlag(value string) (*filter.TimestampFilter, error) {
	if value == "" {
		return nil, fmt.Errorf("timestamp flag: value must not be empty")
	}
	parts := strings.SplitN(value, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("timestamp flag: expected field:layout:min,max — got %q", value)
	}
	field, layout, rangePart := parts[0], parts[1], parts[2]

	bounds := strings.SplitN(rangePart, ",", 2)
	if len(bounds) != 2 {
		return nil, fmt.Errorf("timestamp flag: range must be min,max (either may be empty) — got %q", rangePart)
	}

	var min, max time.Time
	var err error
	if bounds[0] != "" {
		min, err = time.Parse(layout, bounds[0])
		if err != nil {
			return nil, fmt.Errorf("timestamp flag: invalid min %q: %w", bounds[0], err)
		}
	}
	if bounds[1] != "" {
		max, err = time.Parse(layout, bounds[1])
		if err != nil {
			return nil, fmt.Errorf("timestamp flag: invalid max %q: %w", bounds[1], err)
		}
	}

	return filter.NewTimestampFilter(field, layout, min, max)
}
