package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseSamplerFlag parses a sampler flag value like "3" or "every:5".
// Returns the integer N or an error.
func ParseSamplerFlag(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("sampler: empty value")
	}
	if strings.HasPrefix(s, "every:") {
		s = strings.TrimPrefix(s, "every:")
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("sampler: invalid integer %q", s)
	}
	if n < 1 {
		return 0, fmt.Errorf("sampler: n must be >= 1, got %d", n)
	}
	return n, nil
}
