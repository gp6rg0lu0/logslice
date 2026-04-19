package filter

import (
	"testing"

	"github.com/andygrunwald/logslice/internal/parser"
)

func makeHostnameLogLine(field, value string) *parser.LogLine {
	return parser.NewLogLine(map[string]string{field: value})
}

func TestNewHostnameFilter_EmptyField(t *testing.T) {
	_, err := NewHostnameFilter("", []string{"example.com"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewHostnameFilter_NoPatterns(t *testing.T) {
	_, err := NewHostnameFilter("host", []string{})
	if err == nil {
		t.Fatal("expected error for empty patterns")
	}
}

func TestNewHostnameFilter_EmptyPattern(t *testing.T) {
	_, err := NewHostnameFilter("host", []string{"example.com", ""})
	if err == nil {
		t.Fatal("expected error for blank pattern")
	}
}

func TestHostnameFilter_Accessors(t *testing.T) {
	f, err := NewHostnameFilter("host", []string{"*.example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "host" {
		t.Errorf("expected host, got %s", f.Field())
	}
	if len(f.Patterns()) != 1 || f.Patterns()[0] != "*.example.com" {
		t.Errorf("unexpected patterns: %v", f.Patterns())
	}
}

func TestHostnameFilter_NilLine(t *testing.T) {
	f, _ := NewHostnameFilter("host", []string{"example.com"})
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestHostnameFilter_ExactMatch(t *testing.T) {
	f, _ := NewHostnameFilter("host", []string{"example.com"})
	if !f.Match(makeHostnameLogLine("host", "example.com")) {
		t.Error("expected match for exact hostname")
	}
	if f.Match(makeHostnameLogLine("host", "other.com")) {
		t.Error("expected no match")
	}
}

func TestHostnameFilter_WildcardSuffix(t *testing.T) {
	f, _ := NewHostnameFilter("host", []string{"*.example.com"})
	if !f.Match(makeHostnameLogLine("host", "api.example.com")) {
		t.Error("expected match for subdomain")
	}
	if !f.Match(makeHostnameLogLine("host", "example.com")) {
		t.Error("expected match for bare domain")
	}
	if f.Match(makeHostnameLogLine("host", "notexample.com")) {
		t.Error("expected no match")
	}
}

func TestHostnameFilter_WildcardPrefix(t *testing.T) {
	f, _ := NewHostnameFilter("host", []string{"web.*"})
	if !f.Match(makeHostnameLogLine("host", "web.example.com")) {
		t.Error("expected match")
	}
	if !f.Match(makeHostnameLogLine("host", "web")) {
		t.Error("expected match for bare prefix")
	}
	if f.Match(makeHostnameLogLine("host", "webserver.com")) {
		t.Error("expected no match")
	}
}

func TestHostnameFilter_CaseInsensitive(t *testing.T) {
	f, _ := NewHostnameFilter("host", []string{"Example.COM"})
	if !f.Match(makeHostnameLogLine("host", "example.com")) {
		t.Error("expected case-insensitive match")
	}
}
