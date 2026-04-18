package cli

import (
	"fmt"
	"strings"

	"github.com/nicholasgasior/logslice/internal/filter"
)

// ParseJSONPathFlag parses a slice of "path=value" strings into JSONPathFilters.
// Each entry must be in the form "dot.separated.path=value".
func ParseJSONPathFlag(args []string) ([]*filter.JSONPathFilter, error) {
	var filters []*filter.JSONPathFilter
	for _, arg := range args {
		idx := strings.Index(arg, "=")
		if idx <= 0 {
			return nil, fmt.Errorf("jsonpath flag: invalid format %q, expected path=value", arg)
		}
		path := arg[:idx]
		value := arg[idx+1:]
		if value == "" {
			return nil, fmt.Errorf("jsonpath flag: value must not be empty in %q", arg)
		}
		f, err := filter.NewJSONPathFilter(path, value)
		if err != nil {
			return nil, fmt.Errorf("jsonpath flag: %w", err)
		}
		filters = append(filters, f)
	}
	return filters, nil
}
