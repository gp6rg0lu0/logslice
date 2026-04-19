package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// EnvFilter matches log lines where a field value matches one of the given
// environment name patterns (e.g. "prod", "staging", "prod*").
type EnvFilter struct {
	field    string
	envs     []string
	caseFold bool
}

// NewEnvFilter returns a filter that matches when the given field equals one
// of the provided environment strings. Set caseFold to true for
// case-insensitive matching.
func NewEnvFilter(field string, envs []string, caseFold bool) (*EnvFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("envfilter: field must not be empty")
	}
	if len(envs) == 0 {
		return nil, fmt.Errorf("envfilter: at least one environment value required")
	}
	normalized := make([]string, 0, len(envs))
	for _, e := range envs {
		if strings.TrimSpace(e) == "" {
			return nil, fmt.Errorf("envfilter: environment value must not be blank")
		}
		if caseFold {
			normalized = append(normalized, strings.ToLower(e))
		} else {
			normalized = append(normalized, e)
		}
	}
	return &EnvFilter{field: field, envs: normalized, caseFold: caseFold}, nil
}

func (f *EnvFilter) Field() string    { return f.field }
func (f *EnvFilter) Envs() []string   { return f.envs }
func (f *EnvFilter) CaseFold() bool   { return f.caseFold }

// Match returns true when the log line's field value equals one of the
// configured environment strings.
func (f *EnvFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	if f.caseFold {
		val = strings.ToLower(val)
	}
	for _, e := range f.envs {
		if val == e {
			return true
		}
	}
	return false
}
