package cli

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/filter"
)

// ParseTagFlag parses a --tag flag value of the form:
//
//	"field:tag1,tag2,..."
//
// An optional separator override can be supplied as a third colon-delimited
// segment:
//
//	"field:tag1|tag2:|"
//
// When no separator segment is present, "," is used.
func ParseTagFlag(raw string) (*filter.TagFilter, error) {
	if raw == "" {
		return nil, fmt.Errorf("--tag: value must not be empty")
	}

	// Split on ":" with a maximum of 3 parts so that separators containing
	// ":" are not accidentally fragmented.
	parts := strings.SplitN(raw, ":", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("--tag: expected format 'field:tag1,tag2' got %q", raw)
	}

	field := strings.TrimSpace(parts[0])
	if field == "" {
		return nil, fmt.Errorf("--tag: field name must not be empty")
	}

	sep := ","
	if len(parts) == 3 {
		sep = parts[2]
	}

	rawTags := strings.Split(parts[1], sep)
	var tags []string
	for _, t := range rawTags {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}

	return filter.NewTagFilter(field, sep, tags)
}
