package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseLenFlag parses a --len flag value of the form "field:min:max" or "field:min".
// max defaults to -1 (unbounded) when omitted.
// Example: "message:5:100" or "message:10"
func ParseLenFlag(s string) (*filter.LenFilter, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.SplitN(s, ":", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("--len: expected format field:min or field:min:max, got %q", s)
	}
	field := parts[0]
	minVal, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("--len: invalid min value %q: %w", parts[1], err)
	}
	maxVal := -1
	if len(parts) == 3 {
		maxVal, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("--len: invalid max value %q: %w", parts[2], err)
		}
	}
	return filter.NewLenFilter(field, minVal, maxVal)
}
