package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseNotContainsFlag parses a --not-contains flag value of the form "field:substring" or
// "field:substring:i" for case-insensitive matching.
func ParseNotContainsFlag(val string) (*filter.NotContainsFilter, error) {
	if val == "" {
		return nil, nil
	}
	parts := strings.SplitN(val, ":", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("notcontains: expected format field:substring[:i], got %q", val)
	}
	field := parts[0]
	substring := parts[1]
	caseSensitive := true
	if len(parts) == 3 {
		switch strings.ToLower(parts[2]) {
		case "i":
			caseSensitive = false
		case "":
			// default
		default:
			return nil, fmt.Errorf("notcontains: unknown modifier %q, use 'i' for case-insensitive", parts[2])
		}
	}
	return filter.NewNotContainsFilter(field, substring, caseSensitive)
}
