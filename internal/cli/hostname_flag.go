package cli

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/logslice/internal/filter"
)

// ParseHostnameFlag parses the --hostname flag value of the form:
//
//	"field:pattern1,pattern2,..."
//
// and returns a configured HostnameFilter.
func ParseHostnameFlag(value string) (*filter.HostnameFilter, error) {
	if value == "" {
		return nil, nil
	}
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("--hostname: expected format field:pattern1,pattern2 got %q", value)
	}
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return nil, fmt.Errorf("--hostname: field name must not be empty")
	}
	rawPatterns := strings.Split(parts[1], ",")
	var patterns []string
	for _, p := range rawPatterns {
		p = strings.TrimSpace(p)
		if p != "" {
			patterns = append(patterns, p)
		}
	}
	if len(patterns) == 0 {
		return nil, fmt.Errorf("--hostname: at least one pattern is required")
	}
	return filter.NewHostnameFilter(field, patterns)
}
