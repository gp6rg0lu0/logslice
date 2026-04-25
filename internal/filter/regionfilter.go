package filter

import (
	"fmt"
	"strings"

	"github.com/nicholasgasior/logslice/internal/parser"
)

// RegionFilter matches log lines whose field value matches one of the
// specified region identifiers (e.g. "us-east-1", "eu-west-2").
// Matching is case-insensitive and supports optional wildcard suffix
// matching via a trailing "*" (e.g. "us-*" matches any US region).
type RegionFilter struct {
	field    string
	regions  []string
	normalized []string
}

// NewRegionFilter returns a RegionFilter that matches log lines where
// the given field value equals one of the provided region strings.
// Regions may end with "*" to match any region with that prefix.
func NewRegionFilter(field string, regions []string) (*RegionFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("regionfilter: field must not be empty")
	}
	if len(regions) == 0 {
		return nil, fmt.Errorf("regionfilter: at least one region must be specified")
	}
	norm := make([]string, 0, len(regions))
	for _, r := range regions {
		if strings.TrimSpace(r) == "" {
			return nil, fmt.Errorf("regionfilter: region values must not be blank")
		}
		norm = append(norm, strings.ToLower(strings.TrimSpace(r)))
	}
	return &RegionFilter{field: field, regions: regions, normalized: norm}, nil
}

// Field returns the log field this filter inspects.
func (f *RegionFilter) Field() string { return f.field }

// Regions returns the configured region patterns.
func (f *RegionFilter) Regions() []string { return f.regions }

// Match returns true if the log line's field value matches any region pattern.
func (f *RegionFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val := strings.ToLower(strings.TrimSpace(line.Get(f.field)))
	if val == "" {
		return false
	}
	for _, pattern := range f.normalized {
		if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(val, prefix) {
				return true
			}
		} else if val == pattern {
			return true
		}
	}
	return false
}
