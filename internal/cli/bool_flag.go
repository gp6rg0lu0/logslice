package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseBoolFlag parses a repeated --bool flag of the form "field:value"
// where value is "true" or "false". Returns a slice of BoolFilters.
func ParseBoolFlag(values []string) ([]*filter.BoolFilter, error) {
	if len(values) == 0 {
		return nil, nil
	}
	filters := make([]*filter.BoolFilter, 0, len(values))
	for _, v := range values {
		parts := strings.SplitN(v, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("--bool: invalid format %q, expected field:value", v)
		}
		field, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		f, err := filter.NewBoolFilter(field, value)
		if err != nil {
			return nil, fmt.Errorf("--bool: %w", err)
		}
		filters = append(filters, f)
	}
	return filters, nil
}
