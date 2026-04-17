package filter

import (
	"fmt"
	"strconv"

	"github.com/user/logslice/internal/parser"
)

// CompareFilter matches log lines where a numeric field satisfies a comparison.
type CompareFilter struct {
	field string
	op    string
	value float64
}

// NewCompareFilter creates a filter that matches lines where field op value.
// op must be one of: "<", "<=", ">", ">=", "==", "!="
func NewCompareFilter(field, op, raw string) (*CompareFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("comparefilter: field must not be empty")
	}
	switch op {
	case "<", "<=", ">", ">=", "==", "!=":
	default:
		return nil, fmt.Errorf("comparefilter: unsupported operator %q", op)
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return nil, fmt.Errorf("comparefilter: invalid number %q: %w", raw, err)
	}
	return &CompareFilter{field: field, op: op, value: v}, nil
}

func (f *CompareFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	raw, ok := line.Get(f.field)
	if !ok {
		return false
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return false
	}
	switch f.op {
	case "<":
		return v < f.value
	case "<=":
		return v <= f.value
	case ">":
		return v > f.value
	case ">=":
		return v >= f.value
	case "==":
		return v == f.value
	case "!=":
		return v != f.value
	}
	return false
}

func (f *CompareFilter) Field() string  { return f.field }
func (f *CompareFilter) Op() string     { return f.op }
func (f *CompareFilter) Value() float64 { return f.value }
