package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// KeyFilter matches log lines where a field's value exactly matches one of the allowed keys.
type KeyFilter struct {
	field  string
	keys   map[string]struct{}
	raw    string
}

// NewKeyFilter creates a filter that matches lines where field equals one of the comma-separated keys.
func NewKeyFilter(field, keys string) (*KeyFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("keyfilter: field must not be empty")
	}
	if keys == "" {
		return nil, fmt.Errorf("keyfilter: keys must not be empty")
	}
	parts := strings.Split(keys, ",")
	m := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			m[v] = struct{}{}
		}
	}
	if len(m) == 0 {
		return nil, fmt.Errorf("keyfilter: no valid keys provided")
	}
	return &KeyFilter{field: field, keys: m, raw: keys}, nil
}

// Match returns true if the line's field value is one of the allowed keys.
func (f *KeyFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(f.field)
	if !ok {
		return false
	}
	_, found := f.keys[v]
	return found
}

// Field returns the field name being filtered.
func (f *KeyFilter) Field() string { return f.field }

// Keys returns the raw keys string.
func (f *KeyFilter) Keys() string { return f.raw }
