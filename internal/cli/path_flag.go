package cli

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/filter"
)

// ParsePathFlag parses a --path flag value of the form:
//
//	"field:/path1,/path2"         (prefix mode)
//	"field:exact:/path1,/path2"   (exact mode)
//
// Returns a configured *filter.PathFilter or an error.
func ParsePathFlag(value string) (*filter.PathFilter, error) {
	if value == "" {
		return nil, fmt.Errorf("path flag: value must not be empty")
	}

	colon := strings.Index(value, ":")
	if colon < 0 {
		return nil, fmt.Errorf("path flag: expected format field:/p1,/p2 or field:exact:/p1,/p2")
	}

	field := value[:colon]
	rest := value[colon+1:]

	exact := false
	if strings.HasPrefix(rest, "exact:") {
		exact = true
		rest = rest[len("exact:"):]
	}

	parts := strings.Split(rest, ",")
	paths := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			paths = append(paths, p)
		}
	}

	return filter.NewPathFilter(field, paths, exact)
}
