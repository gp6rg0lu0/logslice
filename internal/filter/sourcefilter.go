package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// SourceFilter matches log lines where a field value matches one of the given source identifiers.
type SourceFilter struct {
	field   string
	sources []string
}

// NewSourceFilter creates a filter that matches when the given field equals one of the sources.
func NewSourceFilter(field string, sources []string) (*SourceFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("sourcefilter: field must not be empty")
	}
	var cleaned []string
	for _, s := range sources {
		s = strings.TrimSpace(s)
		if s != "" {
			cleaned = append(cleaned, s)
		}
	}
	if len(cleaned) == 0 {
		return nil, fmt.Errorf("sourcefilter: at least one source value is required")
	}
	return &SourceFilter{field: field, sources: cleaned}, nil
}

// Field returns the field name used for matching.
func (f *SourceFilter) Field() string { return f.field }

// Sources returns the list of accepted source values.
func (f *SourceFilter) Sources() []string { return f.sources }

// Match returns true if the log line's field value matches one of the sources (case-insensitive).
func (f *SourceFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val := strings.ToLower(strings.TrimSpace(line.Get(f.field)))
	for _, s := range f.sources {
		if val == strings.ToLower(s) {
			return true
		}
	}
	return false
}
