package cli

import (
	"testing"
)

func TestParseTimestampFlag_Empty(t *testing.T) {
	_, err := ParseTimestampFlag("")
	if err == nil {
		t.Fatal("expected error for empty value")
	}
}

func TestParseTimestampFlag_MissingParts(t *testing.T) {
	_, err := ParseTimestampFlag("field:layout")
	if err == nil {
		t.Fatal("expected error for missing range part")
	}
}

func TestParseTimestampFlag_MissingRangeComma(t *testing.T) {
	_, err := ParseTimestampFlag("ts:2006-01-02T15:04:05:2024-01-01T00:00:00")
	// SplitN on ":" with n=3 gives [ts, 2006-01-02T15, 04:05:2024-01-01T00:00:00]
	// range part "04:05:2024-01-01T00:00:00" has no comma → error
	if err == nil {
		t.Fatal("expected error when range has no comma")
	}
}

func TestParseTimestampFlag_Valid_Unbounded(t *testing.T) {
	// Use a simple layout to avoid colon-splitting issues
	f, err := ParseTimestampFlag("ts:2006-01-02:,")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "ts" {
		t.Errorf("unexpected field: %s", f.Field())
	}
	if !f.Min().IsZero() || !f.Max().IsZero() {
		t.Error("expected zero min and max for unbounded filter")
	}
}

func TestParseTimestampFlag_Valid_WithBounds(t *testing.T) {
	f, err := ParseTimestampFlag("created:2006-01-02:2024-01-01,2024-12-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "created" {
		t.Errorf("unexpected field: %s", f.Field())
	}
	if f.Min().IsZero() || f.Max().IsZero() {
		t.Error("expected non-zero min and max")
	}
}

func TestParseTimestampFlag_InvalidMin(t *testing.T) {
	_, err := ParseTimestampFlag("ts:2006-01-02:bad-date,2024-12-31")
	if err == nil {
		t.Fatal("expected error for invalid min")
	}
}

func TestParseTimestampFlag_InvalidMax(t *testing.T) {
	_, err := ParseTimestampFlag("ts:2006-01-02:2024-01-01,bad-date")
	if err == nil {
		t.Fatal("expected error for invalid max")
	}
}
