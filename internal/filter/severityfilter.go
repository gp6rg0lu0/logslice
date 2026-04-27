package filter

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// severityRank maps severity label (lowercased) to a numeric rank.
var severityRank = map[string]int{
	"trace":    0,
	"debug":    1,
	"info":     2,
	"notice":   3,
	"warning":  4,
	"warn":     4,
	"error":    5,
	"err":      5,
	"critical": 6,
	"crit":     6,
	"fatal":    7,
	"panic":    8,
}

// SeverityFilter matches log lines whose severity field value falls within
// [minRank, maxRank] (inclusive). Unlike LevelFilter, it supports a configurable
// field name and a finer-grained severity ladder including notice, critical, and panic.
type SeverityFilter struct {
	field   string
	minRank int
	maxRank int
	minLabel string
	maxLabel string
}

// NewSeverityFilter constructs a SeverityFilter for the given field name.
// minSeverity and maxSeverity must be non-empty recognised severity labels.
// Pass "" for minSeverity to default to "trace" and "" for maxSeverity to
// default to "panic".
func NewSeverityFilter(field, minSeverity, maxSeverity string) (*SeverityFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("severityfilter: field must not be empty")
	}
	if minSeverity == "" {
		minSeverity = "trace"
	}
	if maxSeverity == "" {
		maxSeverity = "panic"
	}
	minRank, ok := severityRank[strings.ToLower(minSeverity)]
	if !ok {
		return nil, fmt.Errorf("severityfilter: unknown min severity %q", minSeverity)
	}
	maxRank, ok := severityRank[strings.ToLower(maxSeverity)]
	if !ok {
		return nil, fmt.Errorf("severityfilter: unknown max severity %q", maxSeverity)
	}
	if minRank > maxRank {
		return nil, fmt.Errorf("severityfilter: min severity %q exceeds max severity %q", minSeverity, maxSeverity)
	}
	return &SeverityFilter{
		field:    field,
		minRank:  minRank,
		maxRank:  maxRank,
		minLabel: strings.ToLower(minSeverity),
		maxLabel: strings.ToLower(maxSeverity),
	}, nil
}

// Field returns the log field inspected by this filter.
func (f *SeverityFilter) Field() string { return f.field }

// Min returns the minimum severity label.
func (f *SeverityFilter) Min() string { return f.minLabel }

// Max returns the maximum severity label.
func (f *SeverityFilter) Max() string { return f.maxLabel }

// Match returns true when the line's severity field value falls within [min, max].
func (f *SeverityFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(f.field)
	if !ok {
		return false
	}
	rank, known := severityRank[strings.ToLower(v)]
	if !known {
		return false
	}
	return rank >= f.minRank && rank <= f.maxRank
}
