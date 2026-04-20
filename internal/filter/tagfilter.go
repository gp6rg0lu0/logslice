package filter

import (
	"fmt"
	"strings"

	"github.com/yourorg/logslice/internal/parser"
)

// TagFilter matches log lines where a field contains any of the given tags.
// Tags are compared case-insensitively and the field value is split by a
// configurable separator (default ",").
type TagFilter struct {
	field     string
	tags      map[string]struct{}
	separator string
}

// NewTagFilter creates a TagFilter that matches when the given field contains
// at least one of the provided tags. sep is the delimiter used to split the
// field value; if empty, "," is used.
func NewTagFilter(field, sep string, tags []string) (*TagFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("tagfilter: field must not be empty")
	}
	if len(tags) == 0 {
		return nil, fmt.Errorf("tagfilter: at least one tag is required")
	}
	set := make(map[string]struct{}, len(tags))
	for _, t := range tags {
		t = strings.TrimSpace(t)
		if t == "" {
			return nil, fmt.Errorf("tagfilter: tag must not be blank")
		}
		set[strings.ToLower(t)] = struct{}{}
	}
	if sep == "" {
		sep = ","
	}
	return &TagFilter{field: field, tags: set, separator: sep}, nil
}

// Field returns the field name inspected by this filter.
func (f *TagFilter) Field() string { return f.field }

// Tags returns the set of tags this filter matches against.
func (f *TagFilter) Tags() []string {
	out := make([]string, 0, len(f.tags))
	for t := range f.tags {
		out = append(out, t)
	}
	return out
}

// Match returns true when the log line's field value contains at least one
// of the configured tags after splitting by the separator.
func (f *TagFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	for _, part := range strings.Split(val, f.separator) {
		part = strings.ToLower(strings.TrimSpace(part))
		if _, found := f.tags[part]; found {
			return true
		}
	}
	return false
}
