package filter

import (
	"fmt"
	"time"
)

// TimeFilter filters log entries based on a time range.
type TimeFilter struct {
	From *time.Time
	To   *time.Time
}

// CommonTimeFormats lists the time formats logslice attempts to parse.
var CommonTimeFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// ParseTime attempts to parse a time string using common formats.
func ParseTime(s string) (time.Time, error) {
	for _, layout := range CommonTimeFormats {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse time %q: no matching format found", s)
}

// NewTimeFilter creates a TimeFilter from optional from/to strings.
// Either from or to may be empty, resulting in an open-ended range.
func NewTimeFilter(from, to string) (*TimeFilter, error) {
	tf := &TimeFilter{}

	if from != "" {
		t, err := ParseTime(from)
		if err != nil {
			return nil, fmt.Errorf("invalid --from value: %w", err)
		}
		tf.From = &t
	}

	if to != "" {
		t, err := ParseTime(to)
		if err != nil {
			return nil, fmt.Errorf("invalid --to value: %w", err)
		}
		tf.To = &t
	}

	if tf.From != nil && tf.To != nil && tf.To.Before(*tf.From) {
		return nil, fmt.Errorf("--to (%s) must not be before --from (%s)", to, from)
	}

	return tf, nil
}

// Match returns true if t falls within the filter's time range.
// A nil From or To bound is treated as unbounded.
func (tf *TimeFilter) Match(t time.Time) bool {
	if tf.From != nil && t.Before(*tf.From) {
		return false
	}
	if tf.To != nil && t.After(*tf.To) {
		return false
	}
	return true
}

// Active returns true if at least one bound is set.
func (tf *TimeFilter) Active() bool {
	return tf.From != nil || tf.To != nil
}
