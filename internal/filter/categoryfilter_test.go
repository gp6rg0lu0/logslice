package filter

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
)

func makeCategoryLogLine(field, value string) *parser.LogLine {
	fields := map[string]interface{}{field: value}
	return parser.NewLogLine(fields)
}

func TestNewCategoryFilter_EmptyField(t *testing.T) {
	_, err := NewCategoryFilter("", []string{"auth"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewCategoryFilter_NoCategories(t *testing.T) {
	_, err := NewCategoryFilter("category", nil)
	if err == nil {
		t.Fatal("expected error for nil categories")
	}
}

func TestNewCategoryFilter_BlankCategory(t *testing.T) {
	_, err := NewCategoryFilter("category", []string{"auth", "  "})
	if err == nil {
		t.Fatal("expected error for blank category")
	}
}

func TestCategoryFilter_Accessors(t *testing.T) {
	f, err := NewCategoryFilter("category", []string{"auth", "billing"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "category" {
		t.Errorf("expected field 'category', got %q", f.Field())
	}
	if len(f.Categories()) != 2 {
		t.Errorf("expected 2 categories, got %d", len(f.Categories()))
	}
}

func TestCategoryFilter_NilLine(t *testing.T) {
	f, _ := NewCategoryFilter("category", []string{"auth"})
	if f.Match(nil) {
		t.Error("expected false for nil log line")
	}
}

func TestCategoryFilter_MissingField(t *testing.T) {
	f, _ := NewCategoryFilter("category", []string{"auth"})
	line := makeCategoryLogLine("other", "auth")
	if f.Match(line) {
		t.Error("expected false when field is missing")
	}
}

func TestCategoryFilter_MatchesCaseInsensitive(t *testing.T) {
	f, _ := NewCategoryFilter("category", []string{"Auth", "Billing"})
	cases := []struct {
		val   string
		want  bool
	}{
		{"auth", true},
		{"AUTH", true},
		{"billing", true},
		{"BILLING", true},
		{"payments", false},
		{"", false},
	}
	for _, tc := range cases {
		line := makeCategoryLogLine("category", tc.val)
		if got := f.Match(line); got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.val, got, tc.want)
		}
	}
}

func TestCategoryFilter_DeduplicatesCategories(t *testing.T) {
	f, err := NewCategoryFilter("category", []string{"auth", "Auth", "AUTH"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// All three normalise to the same key; only one unique entry expected.
	if len(f.Categories()) != 1 {
		t.Errorf("expected 1 deduplicated category, got %d", len(f.Categories()))
	}
}
