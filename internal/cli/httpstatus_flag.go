package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// ParseHTTPStatusFlag parses a flag value of the form "field:min-max"
// e.g. "status:400-499" and returns an HTTPStatusFilter.
func ParseHTTPStatusFlag(val string) (*filter.HTTPStatusFilter, error) {
	if val == "" {
		return nil, nil
	}
	colon := strings.IndexByte(val, ':')
	if colon < 1 {
		return nil, fmt.Errorf("httpstatus: expected format field:min-max, got %q", val)
	}
	field := val[:colon]
	rest := val[colon+1:]

	dash := strings.IndexByte(rest, '-')
	if dash < 1 {
		return nil, fmt.Errorf("httpstatus: expected range min-max in %q", rest)
	}
	minStr := rest[:dash]
	maxStr := rest[dash+1:]

	minVal, err := strconv.Atoi(minStr)
	if err != nil {
		return nil, fmt.Errorf("httpstatus: invalid min %q: %w", minStr, err)
	}
	maxVal, err := strconv.Atoi(maxStr)
	if err != nil {
		return nil, fmt.Errorf("httpstatus: invalid max %q: %w", maxStr, err)
	}

	return filter.NewHTTPStatusFilter(field, minVal, maxVal)
}
