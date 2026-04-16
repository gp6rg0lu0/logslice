package filter_test

import (
	"testing"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
)

func makeExistsLogLine(fields map[string]string) *parser.LogLine {
	raw := make(map[string]interface{})
	for k, v := range fields {
		raw[k] = v
	}
	return parser.NewLogLine(raw)
}

func TestNewExistsFilter_EmptyField(t *testing.T) {
	_, err := filter.NewExistsFilter("", true)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewExistsFilter_Valid(t *testing.T) {
	f, err := filter.NewExistsFilter("env", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "env" {
		t.Errorf("expected field 'env', got %q", f.Field())
	}
	if !f.MustExist() {
		t.Error("expected MustExist to be true")
	}
}

func TestExistsFilter_NilLine(t *testing.T) {
	f, _ := filter.NewExistsFilter("env", true)
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestExistsFilter_FieldPresent_MustExistTrue(t *testing.T) {
	f, _ := filter.NewExistsFilter("env", true)
	line := makeExistsLogLine(map[string]string{"env": "prod"})
	if !f.Match(line) {
		t.Error("expected match when field present and mustExist=true")
	}
}

func TestExistsFilter_FieldAbsent_MustExistTrue(t *testing.T) {
	f, _ := filter.NewExistsFilter("env", true)
	line := makeExistsLogLine(map[string]string{"level": "info"})
	if f.Match(line) {
		t.Error("expected no match when field absent and mustExist=true")
	}
}

func TestExistsFilter_FieldAbsent_MustExistFalse(t *testing.T) {
	f, _ := filter.NewExistsFilter("env", false)
	line := makeExistsLogLine(map[string]string{"level": "info"})
	if !f.Match(line) {
		t.Error("expected match when field absent and mustExist=false")
	}
}

func TestExistsFilter_FieldPresent_MustExistFalse(t *testing.T) {
	f, _ := filter.NewExistsFilter("env", false)
	line := makeExistsLogLine(map[string]string{"env": "prod"})
	if f.Match(line) {
		t.Error("expected no match when field present and mustExist=false")
	}
}
