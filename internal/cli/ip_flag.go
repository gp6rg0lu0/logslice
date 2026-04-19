package cli

import (
	"fmt"
	"strings"

	"github.com/nicholasgasior/logslice/internal/filter"
)

// ParseIPFlag parses a slice of "field=cidr" strings and returns a slice of
// IPFilters. Returns an error if any entry is malformed or contains an invalid
// CIDR.
func ParseIPFlag(values []string) ([]*filter.IPFilter, error) {
	var filters []*filter.IPFilter
	for _, v := range values {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("ip flag: expected field=cidr, got %q", v)
		}
		f, err := filter.NewIPFilter(parts[0], parts[1])
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}
	return filters, nil
}
