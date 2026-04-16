package cli

import (
	"fmt"
	"strings"
)

// ParseDedupFlag parses the --dedup flag value.
// The expected format is a single field name, e.g. "request_id".
// Returns the field name or an error if the value is empty or malformed.
func ParseDedupFlag(val string) (string, error) {
	field := strings.TrimSpace(val)
	if field == "" {
		return "", fmt.Errorf("--dedup: field name must not be empty")
	}
	if strings.ContainsAny(field, " \t") {
		return "", fmt.Errorf("--dedup: field name must not contain whitespace")
	}
	return field, nil
}
