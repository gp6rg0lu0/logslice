package filter

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"2024-03-15T10:00:00Z", false},
		{"2024-03-15T10:00:00.123456789Z", false},
		{"2024-03-15T10:00:00", false},
		{"2024-03-15 10:00:00", false},
		{"2024-03-15", false},
		{"not-a-time", true},
		{"", true},
	}

	for _, tc := range cases {
		_, err := ParseTime(tc.input)
		if (err != nil) != tc.wantErr {
			t.Errorf("ParseTime(%q) error = %v, wantErr %v", tc.input, err, tc.wantErr)
		}
	}
}

func TestNewTimeFilter_InvalidRange(t *testing.T) {
	_, err := NewTimeFilter("2024-03-15T12:00:00Z", "2024-03-15T10:00:00Z")
	if err == nil {
		t.Error("expected error when to is before from, got nil")
	}
}

func TestNewTimeFilter_InvalidFormat(t *testing.T) {
	_, err := NewTimeFilter("bad-date", "")
	if err == nil {
		t.Error("expected error for invalid from format, got nil")
	}
}

func TestTimeFilter_Match(t *testing.T) {
	from := time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC)
	to := time.Date(2024, 3, 15, 18, 0, 0, 0, time.UTC)
	tf := &TimeFilter{From: &from, To: &to}

	cases := []struct {
		name  string
		t     time.Time
		want  bool
	}{
		{"before range", time.Date(2024, 3, 15, 7, 59, 59, 0, time.UTC), false},
		{"at from boundary", from, true},
		{"within range", time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC), true},
		{"at to boundary", to, true},
		{"after range", time.Date(2024, 3, 15, 18, 0, 1, 0, time.UTC), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tf.Match(tc.t); got != tc.want {
				t.Errorf("Match() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestTimeFilter_Active(t *testing.T) {
	empty := &TimeFilter{}
	if empty.Active() {
		t.Error("expected Active() = false for empty filter")
	}

	tf, _ := NewTimeFilter("2024-01-01", "")
	if !tf.Active() {
		t.Error("expected Active() = true when From is set")
	}
}
