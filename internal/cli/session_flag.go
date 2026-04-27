package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseSessionFlag parses the --session flag value into a SessionFilter.
//
// Format: "field:value1,value2,..." or "field:prefix:value1,value2,..."
//
// Examples:
//
//	--session session_id:abc123,def456          (exact match)
//	--session session_id:prefix:sess-a,sess-b   (prefix match)
func ParseSessionFlag(s string) (*filter.SessionFilter, error) {
	if s == "" {
		return nil, nil
	}

	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("--session: expected format field:value[,...] or field:prefix:value[,...], got %q", s)
	}

	field := strings.TrimSpace(parts[0])
	if field == "" {
		return nil, fmt.Errorf("--session: field name must not be empty")
	}

	rest := parts[1]
	prefixMode := false

	if strings.HasPrefix(rest, "prefix:") {
		prefixMode = true
		rest = strings.TrimPrefix(rest, "prefix:")
	}

	var values []string
	for _, v := range strings.Split(rest, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			values = append(values, v)
		}
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("--session: at least one session value is required")
	}

	f, err := filter.NewSessionFilter(field, values, prefixMode)
	if err != nil {
		return nil, fmt.Errorf("--session: %w", err)
	}
	return f, nil
}
