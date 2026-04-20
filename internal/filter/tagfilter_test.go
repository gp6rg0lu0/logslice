package filter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
)

func makeTagLogLine(fields map[string]string) *parser.LogLine {
	m := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		m[k] = v
	}
	return parser.NewLogLine(m)
}

func TestNewTagFilter_EmptyField(t *testing.T) {
	_, err := filter.NewTagFilter("", ",", []string{"prod"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewTagFilter_NoTags(t *testing.T) {
	_, err := filter.NewTagFilter("tags", ",", nil)
	if err == nil {
		t.Fatal("expected error for empty tag list")
	}
}

func TestNewTagFilter_BlankTag(t *testing.T) {
	_, err := filter.NewTagFilter("tags", ",", []string{"prod", ""})
	if err == nil {
		t.Fatal("expected error for blank tag")
	}
}

func TestTagFilter_Accessors(t *testing.T) {
	f, err := filter.NewTagFilter("tags", ",", []string{"prod", "us-east"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "tags" {
		t.Errorf("expected field 'tags', got %q", f.Field())
	}
	if len(f.Tags()) != 2 {
		t.Errorf("expected 2 tags, got %d", len(f.Tags()))
	}
}

func TestTagFilter_NilLine(t *testing.T) {
	f, _ := filter.NewTagFilter("tags", ",", []string{"prod"})
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestTagFilter_MissingField(t *testing.T) {
	f, _ := filter.NewTagFilter("tags", ",", []string{"prod"})
	line := makeTagLogLine(map[string]string{"level": "info"})
	if f.Match(line) {
		t.Error("expected false when field is missing")
	}
}

func TestTagFilter_MatchesSingleTag(t *testing.T) {
	f, _ := filter.NewTagFilter("tags", ",", []string{"prod"})
	line := makeTagLogLine(map[string]string{"tags": "prod,staging"})
	if !f.Match(line) {
		t.Error("expected true for matching tag")
	}
}

func TestTagFilter_CaseInsensitive(t *testing.T) {
	f, _ := filter.NewTagFilter("tags", ",", []string{"PROD"})
	line := makeTagLogLine(map[string]string{"tags": "prod,staging"})
	if !f.Match(line) {
		t.Error("expected case-insensitive match")
	}
}

func TestTagFilter_NoMatch(t *testing.T) {
	f, _ := filter.NewTagFilter("tags", ",", []string{"canary"})
	line := makeTagLogLine(map[string]string{"tags": "prod,staging"})
	if f.Match(line) {
		t.Error("expected false when no tag matches")
	}
}

func TestTagFilter_CustomSeparator(t *testing.T) {
	f, _ := filter.NewTagFilter("tags", "|", []string{"us-east"})
	line := makeTagLogLine(map[string]string{"tags": "prod|us-east|v2"})
	if !f.Match(line) {
		t.Error("expected true with custom separator")
	}
}

func TestTagFilter_DefaultSeparatorUsedWhenEmpty(t *testing.T) {
	f, err := filter.NewTagFilter("tags", "", []string{"prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeTagLogLine(map[string]string{"tags": "prod,staging"})
	if !f.Match(line) {
		t.Error("expected true using default comma separator")
	}
}
