package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseEnvFlag parses a --env flag value of the form:
//
//	"field:val1,val2,..."
//
// An optional "+i" suffix on the field name enables case-insensitive matching:
//
//	"field+i:val1,val2"
func ParseEnvFlag(s string) (*filter.EnvFilter, error) {
	if s == "" {
		return nil, fmt.Errorf("env flag: value must not be empty")
	}
	colon := strings.Index(s, ":")
	if colon < 0 {
		return nil, fmt.Errorf("env flag: expected format field:val1,val2 — got %q", s)
	}
	fieldPart := s[:colon]
	valPart := s[colon+1:]

	caseFold := false
	if strings.HasSuffix(fieldPart, "+i") {
		caseFold = true
		fieldPart = strings.TrimSuffix(fieldPart, "+i")
	}

	if fieldPart == "" {
		return nil, fmt.Errorf("env flag: field name must not be empty")
	}
	if valPart == "" {
		return nil, fmt.Errorf("env flag: at least one value required after ':'")
	}

	envs := strings.Split(valPart, ",")
	return filter.NewEnvFilter(fieldPart, envs, caseFold)
}
