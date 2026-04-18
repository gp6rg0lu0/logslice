package cli

import (
	"fmt"
	"strings"

	"github.com/nicholasgasior/logslice/internal/filter"
)

// ParseTypeFlag parses a --type flag value of the form "field:typename".
// Example: --type meta:object
func ParseTypeFlag(s string) (*filter.TypeFilter, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("--type: expected format field:typename, got %q", s)
	}
	f, err := filter.NewTypeFilter(parts[0], parts[1])
	if err != nil {
		return nil, fmt.Errorf("--type: %w", err)
	}
	return f, nil
}
