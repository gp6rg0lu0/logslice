package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseKeyFlag parses a --key flag value of the form "field=val1,val2,..."
// and returns a KeyFilter.
func ParseKeyFlag(value string) (*filter.KeyFilter, error) {
	if value == "" {
		return nil, nil
	}
	idx := strings.IndexByte(value, '=')
	if idx < 1 {
		return nil, fmt.Errorf("--key: expected format field=val1,val2 but got %q", value)
	}
	field := strings.TrimSpace(value[:idx])
	keys := strings.TrimSpace(value[idx+1:])
	if field == "" {
		return nil, fmt.Errorf("--key: field name must not be empty")
	}
	if keys == "" {
		return nil, fmt.Errorf("--key: keys must not be empty")
	}
	return filter.NewKeyFilter(field, keys)
}
