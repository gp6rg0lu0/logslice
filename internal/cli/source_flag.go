package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseSourceFlag parses a --source flag value of the form "field:val1,val2,..."
// and returns a SourceFilter.
func ParseSourceFlag(value string) (*filter.SourceFilter, error) {
	if value == "" {
		return nil, nil
	}
	idx := strings.Index(value, ":")
	if idx < 1 {
		return nil, fmt.Errorf("--source: expected format field:val1,val2 got %q", value)
	}
	field := value[:idx]
	raw := value[idx+1:]
	if raw == "" {
		return nil, fmt.Errorf("--source: no values provided for field %q", field)
	}
	sources := strings.Split(raw, ",")
	return filter.NewSourceFilter(field, sources)
}
