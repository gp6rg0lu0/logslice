package cli

import (
	"fmt"
	"strings"

	"github.com/nicholasgasior/logslice/internal/filter"
)

// ParseRegionFlag parses a --region flag value of the form:
//
//	"field:region1,region2,..."
//
// Regions may include a trailing "*" wildcard (e.g. "us-*").
// Returns a RegionFilter or an error if the format is invalid.
func ParseRegionFlag(value string) (*filter.RegionFilter, error) {
	if value == "" {
		return nil, fmt.Errorf("region flag: value must not be empty")
	}
	idx := strings.Index(value, ":")
	if idx < 0 {
		return nil, fmt.Errorf("region flag: expected format field:region1,region2 — missing colon")
	}
	field := strings.TrimSpace(value[:idx])
	if field == "" {
		return nil, fmt.Errorf("region flag: field name must not be empty")
	}
	raw := value[idx+1:]
	parts := strings.Split(raw, ",")
	regions := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			regions = append(regions, t)
		}
	}
	if len(regions) == 0 {
		return nil, fmt.Errorf("region flag: at least one region must be specified")
	}
	return filter.NewRegionFilter(field, regions)
}
