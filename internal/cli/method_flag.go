package cli

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/filter"
)

// ParseMethodFlag parses a --method flag value of the form:
//
//	"field:METHOD1,METHOD2,..."
//
// For example: "http_method:GET,POST" creates a MethodFilter that matches
// log lines where the http_method field is GET or POST.
func ParseMethodFlag(value string) (*filter.MethodFilter, error) {
	if value == "" {
		return nil, nil
	}

	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("--method: expected format field:METHOD1,METHOD2 got %q", value)
	}

	field := strings.TrimSpace(parts[0])
	if field == "" {
		return nil, fmt.Errorf("--method: field name must not be empty")
	}

	rawMethods := strings.Split(parts[1], ",")
	methods := make([]string, 0, len(rawMethods))
	for _, m := range rawMethods {
		m = strings.TrimSpace(m)
		if m != "" {
			methods = append(methods, m)
		}
	}

	if len(methods) == 0 {
		return nil, fmt.Errorf("--method: at least one HTTP method must be specified")
	}

	return filter.NewMethodFilter(field, methods)
}
