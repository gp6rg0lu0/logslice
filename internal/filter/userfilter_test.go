package filter

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func makeUserLogLine(field, value string) *parser.LogLine {
	data := map[string]string{field: value}
	return parser.NewLogLine(data)
}

func TestNewUserFilter_EmptyField(t *testing.T) {
	_, err := NewUserFilter("", []string{"alice"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewUserFilter_NoValues(t *testing.T) {
	_, err := NewUserFilter("user", []string{})
	if err == nil {
		t.Fatal("expected error for empty user list")
	}
}

func TestNewUserFilter_BlankValue(t *testing.T) {
	_, err := NewUserFilter("user", []string{"alice", "  "})
	if err == nil {
		t.Fatal("expected error for blank user value")
	}
}

func TestNewUserFilter_Valid(t *testing.T) {
	f, err := NewUserFilter("user", []string{"alice", "bob"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "user" {
		t.Errorf("expected field 'user', got %q", f.Field())
	}
	if len(f.Users()) != 2 {
		t.Errorf("expected 2 users, got %d", len(f.Users()))
	}
}

func TestUserFilter_NilLine(t *testing.T) {
	f, _ := NewUserFilter("user", []string{"alice"})
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestUserFilter_MissingField(t *testing.T) {
	f, _ := NewUserFilter("user", []string{"alice"})
	line := makeUserLogLine("actor", "alice")
	if f.Match(line) {
		t.Error("expected false when field is missing")
	}
}

func TestUserFilter_MatchesCaseInsensitive(t *testing.T) {
	f, _ := NewUserFilter("user", []string{"Alice"})
	line := makeUserLogLine("user", "alice")
	if !f.Match(line) {
		t.Error("expected match for case-insensitive value")
	}
}

func TestUserFilter_NoMatch(t *testing.T) {
	f, _ := NewUserFilter("user", []string{"alice", "bob"})
	line := makeUserLogLine("user", "charlie")
	if f.Match(line) {
		t.Error("expected no match for unknown user")
	}
}

func TestUserFilter_MatchesMultiple(t *testing.T) {
	f, _ := NewUserFilter("user", []string{"alice", "bob"})
	for _, name := range []string{"alice", "bob", "BOB", "Alice"} {
		line := makeUserLogLine("user", name)
		if !f.Match(line) {
			t.Errorf("expected match for user %q", name)
		}
	}
}
