package filter

import (
	"fmt"

	"github.com/user/logslice/internal/parser"
)

// BoolFilter matches log lines where a field equals a boolean value.
type BoolFilter struct {
	field string
	want  bool
}

// NewBoolFilter returns a filter that matches lines where field equals want.
// The value string must be "true" or "false".
func NewBoolFilter(field, value string) (*BoolFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("boolfilter: field must not be empty")
	}
	switch value {
	case "true":
		return &BoolFilter{field: field, want: true}, nil
	case "false":
		return &BoolFilter{field: field, want: false}, nil
	default:
		return nil, fmt.Errorf("boolfilter: value must be \"true\" or \"false\", got %q", value)
	}
}

// Field returns the field name.
func (f *BoolFilter) Field() string { return f.field }

// Want returns the expected boolean value.
func (f *BoolFilter) Want() bool { return f.want }

// Match returns true if the field value equals the expected boolean.
func (f *BoolFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(f.field)
	if !ok {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val == f.want
	case string:
		return (val == "true") == f.want
	default:
		return false
	}
}
