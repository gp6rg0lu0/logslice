package cli

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// ParseUserFlag parses a --user flag value of the form "field:user1,user2,..."
// and returns a configured UserFilter.
//
// Example:
//
//	--user "username:alice,bob"
func ParseUserFlag(value string) (*filter.UserFilter, error) {
	if value == "" {
		return nil, fmt.Errorf("user flag: value must not be empty")
	}
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("user flag: expected format field:user1,user2 — got %q", value)
	}
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return nil, fmt.Errorf("user flag: field name must not be empty")
	}
	raw := strings.Split(parts[1], ",")
	users := make([]string, 0, len(raw))
	for _, u := range raw {
		trimmed := strings.TrimSpace(u)
		if trimmed != "" {
			users = append(users, trimmed)
		}
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("user flag: at least one user value is required")
	}
	return filter.NewUserFilter(field, users)
}
