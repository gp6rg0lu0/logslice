package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseTraceFlag parses a --trace flag value of the form:
//   field:value1,value2[,exact]
// If the last segment is the literal "exact", exact matching is enabled.
func ParseTraceFlag(s string) (*filter.TraceFilter, error) {
	if s == "" {
		return nil, fmt.Errorf("trace flag must not be empty")
	}
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("trace flag must be in form field:value1,value2[,exact]")
	}
	field := strings.TrimSpace(parts[0])
	rest := parts[1]

	segments := strings.Split(rest, ",")
	exact := false
	if len(segments) > 0 && strings.TrimSpace(segments[len(segments)-1]) == "exact" {
		exact = true
		segments = segments[:len(segments)-1]
	}
	var values []string
	for _, seg := range segments {
		v := strings.TrimSpace(seg)
		if v != "" {
			values = append(values, v)
		}
	}
	return filter.NewTraceFilter(field, values, exact)
}
