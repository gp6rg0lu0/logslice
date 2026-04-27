package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// SessionFilter matches log lines where a session field contains one of the
// specified session IDs, with optional prefix matching.
type SessionFilter struct {
	field    string
	values   []string
	prefix   bool
}

// NewSessionFilter creates a SessionFilter for the given field and session IDs.
// If prefix is true, values are treated as prefixes rather than exact matches.
func NewSessionFilter(field string, values []string, prefix bool) (*SessionFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("sessionfilter: field must not be empty")
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("sessionfilter: at least one session value is required")
	}
	var clean []string
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			return nil, fmt.Errorf("sessionfilter: session values must not be blank")
		}
		clean = append(clean, v)
	}
	return &SessionFilter{field: field, values: clean, prefix: prefix}, nil
}

// Field returns the log field inspected by the filter.
func (f *SessionFilter) Field() string { return f.field }

// Values returns the session IDs (or prefixes) used for matching.
func (f *SessionFilter) Values() []string { return f.values }

// Prefix reports whether prefix matching is enabled.
func (f *SessionFilter) Prefix() bool { return f.prefix }

// Match returns true if the line's session field equals (or starts with, when
// prefix mode is active) one of the configured values.
func (f *SessionFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v := line.Get(f.field)
	if v == "" {
		return false
	}
	for _, s := range f.values {
		if f.prefix {
			if strings.HasPrefix(v, s) {
				return true
			}
		} else {
			if v == s {
				return true
			}
		}
	}
	return false
}
