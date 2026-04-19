package filter

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

const tsLayout = "2006-01-02T15:04:05"

func makeTimestampLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]interface{}{field: value})
}

func TestNewTimestampFilter_EmptyField(t *testing.T) {
	_, err := NewTimestampFilter("", tsLayout, time.Time{}, time.Time{})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewTimestampFilter_EmptyLayout(t *testing.T) {
	_, err := NewTimestampFilter("ts", "", time.Time{}, time.Time{})
	if err == nil {
		t.Fatal("expected error for empty layout")
	}
}

func TestNewTimestampFilter_MaxBeforeMin(t *testing.T) {
	min := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	max := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := NewTimestampFilter("ts", tsLayout, min, max)
	if err == nil {
		t.Fatal("expected error when max before min")
	}
}

func TestTimestampFilter_Accessors(t *testing.T) {
	min := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	f, _ := NewTimestampFilter("ts", tsLayout, min, time.Time{})
	if f.Field() != "ts" {
		t.Errorf("unexpected field: %s", f.Field())
	}
	if f.Layout() != tsLayout {
		t.Errorf("unexpected layout: %s", f.Layout())
	}
	if !f.Min().Equal(min) {
		t.Errorf("unexpected min")
	}
}

func TestTimestampFilter_NilLine(t *testing.T) {
	f, _ := NewTimestampFilter("ts", tsLayout, time.Time{}, time.Time{})
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestTimestampFilter_MatchWithinRange(t *testing.T) {
	min := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	max := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	f, _ := NewTimestampFilter("ts", tsLayout, min, max)
	line := makeTimestampLogLine("ts", "2024-06-15T12:00:00")
	if !f.Match(line) {
		t.Error("expected match for timestamp within range")
	}
}

func TestTimestampFilter_MatchOutsideRange(t *testing.T) {
	min := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	max := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	f, _ := NewTimestampFilter("ts", tsLayout, min, max)
	line := makeTimestampLogLine("ts", "2023-06-15T12:00:00")
	if f.Match(line) {
		t.Error("expected no match for timestamp outside range")
	}
}

func TestTimestampFilter_InvalidTimestamp(t *testing.T) {
	f, _ := NewTimestampFilter("ts", tsLayout, time.Time{}, time.Time{})
	line := makeTimestampLogLine("ts", "not-a-time")
	if f.Match(line) {
		t.Error("expected no match for unparseable timestamp")
	}
}
