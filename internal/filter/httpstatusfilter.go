package filter

import (
	"fmt"
	"strconv"

	"github.com/yourorg/logslice/internal/parser"
)

// HTTPStatusFilter matches log lines where a field's HTTP status code falls
// within a given range (e.g. 400–499 for client errors).
type HTTPStatusFilter struct {
	field  string
	minVal int
	maxVal int
}

// NewHTTPStatusFilter creates a filter that matches when the integer value of
// field is between minVal and maxVal (inclusive).
func NewHTTPStatusFilter(field string, minVal, maxVal int) (*HTTPStatusFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("httpstatus: field name must not be empty")
	}
	if minVal < 100 || maxVal > 599 {
		return nil, fmt.Errorf("httpstatus: status codes must be between 100 and 599")
	}
	if minVal > maxVal {
		return nil, fmt.Errorf("httpstatus: min %d exceeds max %d", minVal, maxVal)
	}
	return &HTTPStatusFilter{field: field, minVal: minVal, maxVal: maxVal}, nil
}

func (f *HTTPStatusFilter) Field() string  { return f.field }
func (f *HTTPStatusFilter) Min() int       { return f.minVal }
func (f *HTTPStatusFilter) Max() int       { return f.maxVal }

func (f *HTTPStatusFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(f.field)
	if !ok {
		return false
	}
	code, err := strconv.Atoi(v)
	if err != nil {
		return false
	}
	return code >= f.minVal && code <= f.maxVal
}
