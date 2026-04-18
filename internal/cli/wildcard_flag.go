package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseWildcardFlag parses a repeated --wildcard flag of the form "field=pattern"
// and returns a slice of WildcardFilters.
func ParseWildcardFlag(values []string) ([]*filter.WildcardFilter, error) {
	if len(values) == 0 {
		return nil, nil
	}
	filters := make([]*filter.WildcardFilter, 0, len(values))
	for _, v := range values {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("wildcard flag: invalid format %q, expected field=pattern", v)
		}
		field, pattern := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		f, err := filter.NewWildcardFilter(field, pattern)
		if err != nil {
			return nil, fmt.Errorf("wildcard flag: %w", err)
		}
		filters = append(filters, f)
	}
	return filters, nil
}
