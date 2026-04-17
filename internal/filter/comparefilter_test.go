package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeCompareLogLine(fields map[string]string) *parser.LogLine {
	return parser.NewLogLine(fields)
}

func TestNewCompareFilter_EmptyField(t *testing.T) {
	_, err := NewCompareFilter("", ">", "5")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewCompareFilter_BadOp(t *testing.T) {
	_, err := NewCompareFilter("latency", "??", "5")
	if err == nil {
		t.Fatal("expected error for bad operator")
	}
}

func TestNewCompareFilter_BadNumber(t *testing.T) {
	_, err := NewCompareFilter("latency", ">", "abc")
	if err == nil {
		t.Fatal("expected error for non-numeric value")
	}
}

func TestCompareFilter_Accessors(t *testing.T) {
	f, err := NewCompareFilter("latency", ">=", "100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "latency" {
		t.Errorf("expected field latency, got %s", f.Field())
	}
	if f.Op() != ">=" {
		t.Errorf("expected op >=, got %s", f.Op())
	}
	if f.Value() != 100 {
		t.Errorf("expected value 100, got %f", f.Value())
	}
}

func TestCompareFilter_NilLine(t *testing.T) {
	f, _ := NewCompareFilter("latency", ">", "0")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestCompareFilter_MissingField(t *testing.T) {
	f, _ := NewCompareFilter("latency", ">", "0")
	line := makeCompareLogLine(map[string]string{"msg": "hello"})
	if f.Match(line) {
		t.Error("expected false for missing field")
	}
}

func TestCompareFilter_Match(t *testing.T) {
	tests := []struct {
		op    string
		val   string
		field string
		want  bool
	}{
		{">", "50", "100", true},
		{">", "100", "100", false},
		{">=", "100", "100", true},
		{"<", "200", "100", true},
		{"<=", "100", "100", true},
		{"==", "100", "100", true},
		{"==", "99", "100", false},
		{"!=", "99", "100", true},
		{"!=", "100", "100", false},
	}
	for _, tt := range tests {
		f, err := NewCompareFilter("latency", tt.op, tt.val)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		line := makeCompareLogLine(map[string]string{"latency": tt.field})
		if got := f.Match(line); got != tt.want {
			t.Errorf("op=%s val=%s field=%s: expected %v got %v", tt.op, tt.val, tt.field, tt.want, got)
		}
	}
}
