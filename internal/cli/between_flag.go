package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseBetweenFlag parses a slice of "field:min:max" strings into BetweenFilters.
// Each entry must have exactly three colon-separated parts.
func ParseBetweenFlag(args []string) ([]*filter.BetweenFilter, error) {
	var filters []*filter.BetweenFilter
	for _, arg := range args {
		parts := strings.SplitN(arg, ":", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("between: invalid format %q, expected field:min:max", arg)
		}
		field := parts[0]
		min, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("between: invalid min %q in %q", parts[1], arg)
		}
		max, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return nil, fmt.Errorf("between: invalid max %q in %q", parts[2], arg)
		}
		f, err := filter.NewBetweenFilter(field, min, max)
		if err != nil {
			return nil, fmt.Errorf("between: %w", err)
		}
		filters = append(filters, f)
	}
	return filters, nil
}
