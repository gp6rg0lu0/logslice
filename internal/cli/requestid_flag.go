package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseRequestIDFlag parses the --request-id flag value into a RequestIDFilter.
//
// Format: "field:mode:id1,id2,..."
//   - field: the log field to inspect (e.g. "request_id")
//   - mode:  "exact" or "prefix"
//   - ids:   comma-separated list of request ID values to match
//
// Example: "request_id:prefix:req-abc,req-def"
func ParseRequestIDFlag(value string) (*filter.RequestIDFilter, error) {
	if value == "" {
		return nil, nil
	}
	parts := strings.SplitN(value, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("--request-id: expected format field:mode:id1,id2,... got %q", value)
	}
	field := strings.TrimSpace(parts[0])
	mode := strings.TrimSpace(parts[1])
	rawIDs := strings.TrimSpace(parts[2])

	if field == "" {
		return nil, fmt.Errorf("--request-id: field must not be empty")
	}
	if rawIDs == "" {
		return nil, fmt.Errorf("--request-id: at least one request ID must be provided")
	}

	var exact bool
	switch mode {
	case "exact":
		exact = true
	case "prefix":
		exact = false
	default:
		return nil, fmt.Errorf("--request-id: mode must be \"exact\" or \"prefix\", got %q", mode)
	}

	ids := strings.Split(rawIDs, ",")
	return filter.NewRequestIDFilter(field, ids, exact)
}
