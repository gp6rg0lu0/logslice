package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseCompareFlag parses a compare flag value of the form "field:op:number".
// Example: "latency:>:100" or "status:==:200"
func ParseCompareFlag(s string) (*filter.CompareFilter, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.SplitN(s, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("compare flag: expected field:op:value, got %q", s)
	}
	field, op, value := parts[0], parts[1], parts[2]
	f, err := filter.NewCompareFilter(field, op, value)
	if err != nil {
		return nil, fmt.Errorf("compare flag: %w", err)
	}
	return f, nil
}
