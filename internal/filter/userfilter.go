package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// UserFilter matches log lines where a specified field's value matches
// one of the provided user identifiers (exact, case-insensitive match).
type UserFilter struct {
	field string
	users map[string]struct{}
}

// NewUserFilter creates a UserFilter that matches when the given field
// contains any of the provided user values. Field must be non-empty and
// at least one non-blank user value must be supplied.
func NewUserFilter(field string, users []string) (*UserFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("userfilter: field must not be empty")
	}
	set := make(map[string]struct{}, len(users))
	for _, u := range users {
		trimmed := strings.TrimSpace(u)
		if trimmed == "" {
			return nil, fmt.Errorf("userfilter: user value must not be blank")
		}
		set[strings.ToLower(trimmed)] = struct{}{}
	}
	if len(set) == 0 {
		return nil, fmt.Errorf("userfilter: at least one user value is required")
	}
	return &UserFilter{field: field, users: set}, nil
}

// Field returns the log field inspected by this filter.
func (f *UserFilter) Field() string { return f.field }

// Users returns a sorted slice of the user values this filter matches.
func (f *UserFilter) Users() []string {
	out := make([]string, 0, len(f.users))
	for u := range f.users {
		out = append(out, u)
	}
	return out
}

// Match returns true when the log line's field value matches one of the
// configured user identifiers (case-insensitive).
func (f *UserFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	_, found := f.users[strings.ToLower(strings.TrimSpace(val))]
	return found
}
