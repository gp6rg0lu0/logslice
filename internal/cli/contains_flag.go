package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseContainsFlag parses a --contains flag value of the form "field:substring"
// or "field:substring:fold" where fold=true enables case-insensitive matching.
func ParseContainsFlag(val string) (*filter.ContainsFilter, error) {
	if val == "" {
		return nil, nil
	}
	parts := strings.SplitN(val, ":", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("--contains: expected field:substring[:fold], got %q", val)
	}
	field := parts[0]
	substring := parts[1]
	caseFold := false
	if len(parts) == 3 {
		switch strings.ToLower(parts[2]) {
		case "true", "1", "fold":
			caseFold = true
		case "false", "0", "":
			caseFold = false
		default:
			return nil, fmt.Errorf("--contains: invalid fold value %q, use true/false", parts[2])
		}
	}
	return filter.NewContainsFilter(field, substring, caseFold)
}
