package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseMultiValueFlag parses a flag value of the form "field:val1,val2,..." into a MultiValueFilter.
// Example: --in "level:info,warn,error"
func ParseMultiValueFlag(s string) (*filter.MultiValueFilter, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("--in: expected format field:val1,val2 got %q", s)
	}
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return nil, fmt.Errorf("--in: field name must not be empty")
	}
	rawValues := strings.Split(parts[1], ",")
	var values []string
	for _, v := range rawValues {
		v = strings.TrimSpace(v)
		if v != "" {
			values = append(values, v)
		}
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("--in: at least one value required in %q", s)
	}
	return filter.NewMultiValueFilter(field, values)
}
