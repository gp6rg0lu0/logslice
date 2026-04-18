package filter_test

import (
	"testing"

	"github.com/nicholasgasior/logslice/internal/filter"
	"github.com/nicholasgasior/logslice/internal/parser"
)

func makeTypeLogLine(fields map[string]interface{}) *parser.LogLine {
	return parser.NewLogLineRaw(fields)
}

func TestNewTypeFilter_EmptyField(t *testing.T) {
	_, err := filter.NewTypeFilter("", "string")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewTypeFilter_InvalidType(t *testing.T) {
	_, err := filter.NewTypeFilter("level", "integer")
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestNewTypeFilter_Valid(t *testing.T) {
	f, err := filter.NewTypeFilter("count", "number")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "count" || f.WantType() != "number" {
		t.Errorf("accessors returned wrong values")
	}
}

func TestTypeFilter_NilLine(t *testing.T) {
	f, _ := filter.NewTypeFilter("x", "string")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestTypeFilter_Match(t *testing.T) {
	tests := []struct {
		name     string
		fields   map[string]interface{}
		field    string
		wantType string
		expect   bool
	}{
		{"string match", map[string]interface{}{"msg": "hello"}, "msg", "string", true},
		{"number match", map[string]interface{}{"count": float64(3)}, "count", "number", true},
		{"bool match", map[string]interface{}{"ok": true}, "ok", "bool", true},
		{"null match", map[string]interface{}{"x": nil}, "x", "null", true},
		{"object match", map[string]interface{}{"meta": map[string]interface{}{"k": "v"}}, "meta", "object", true},
		{"array match", map[string]interface{}{"tags": []interface{}{"a"}}, "tags", "array", true},
		{"type mismatch", map[string]interface{}{"count": float64(1)}, "count", "string", false},
		{"missing field null", map[string]interface{}{}, "x", "null", true},
		{"missing field string", map[string]interface{}{}, "x", "string", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := filter.NewTypeFilter(tt.field, tt.wantType)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			line := makeTypeLogLine(tt.fields)
			if got := f.Match(line); got != tt.expect {
				t.Errorf("Match() = %v, want %v", got, tt.expect)
			}
		})
	}
}
