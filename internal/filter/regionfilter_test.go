package filter

import (
	"testing"

	"github.com/nicholasgasior/logslice/internal/parser"
)

func makeRegionLogLine(field, value string) *parser.LogLine {
	data := map[string]interface{}{field: value}
	return parser.NewLogLine(data)
}

func TestNewRegionFilter_EmptyField(t *testing.T) {
	_, err := NewRegionFilter("", []string{"us-east-1"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewRegionFilter_NoRegions(t *testing.T) {
	_, err := NewRegionFilter("region", []string{})
	if err == nil {
		t.Fatal("expected error for empty regions slice")
	}
}

func TestNewRegionFilter_BlankRegion(t *testing.T) {
	_, err := NewRegionFilter("region", []string{"us-east-1", "  "})
	if err == nil {
		t.Fatal("expected error for blank region value")
	}
}

func TestRegionFilter_Accessors(t *testing.T) {
	f, err := NewRegionFilter("region", []string{"us-east-1", "eu-west-2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "region" {
		t.Errorf("expected field 'region', got %q", f.Field())
	}
	if len(f.Regions()) != 2 {
		t.Errorf("expected 2 regions, got %d", len(f.Regions()))
	}
}

func TestRegionFilter_NilLine(t *testing.T) {
	f, _ := NewRegionFilter("region", []string{"us-east-1"})
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestRegionFilter_ExactMatch(t *testing.T) {
	f, _ := NewRegionFilter("region", []string{"us-east-1", "eu-west-2"})
	if !f.Match(makeRegionLogLine("region", "us-east-1")) {
		t.Error("expected match for us-east-1")
	}
	if !f.Match(makeRegionLogLine("region", "EU-WEST-2")) {
		t.Error("expected case-insensitive match for EU-WEST-2")
	}
	if f.Match(makeRegionLogLine("region", "ap-southeast-1")) {
		t.Error("expected no match for ap-southeast-1")
	}
}

func TestRegionFilter_WildcardMatch(t *testing.T) {
	f, _ := NewRegionFilter("region", []string{"us-*"})
	if !f.Match(makeRegionLogLine("region", "us-east-1")) {
		t.Error("expected wildcard match for us-east-1")
	}
	if !f.Match(makeRegionLogLine("region", "us-west-2")) {
		t.Error("expected wildcard match for us-west-2")
	}
	if f.Match(makeRegionLogLine("region", "eu-west-1")) {
		t.Error("expected no match for eu-west-1 with us-* pattern")
	}
}

func TestRegionFilter_MissingField(t *testing.T) {
	f, _ := NewRegionFilter("region", []string{"us-east-1"})
	if f.Match(makeRegionLogLine("zone", "us-east-1")) {
		t.Error("expected no match when field is missing")
	}
}
