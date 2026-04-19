package filter

import (
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func makeHashLogLine(fields map[string]string) *parser.LogLine {
	m := make(map[string]interface{})
	for k, v := range fields {
		m[k] = v
	}
	return parser.NewLogLine(m)
}

func TestNewHashFilter_NoFields(t *testing.T) {
	_, err := NewHashFilter([]string{})
	if err == nil {
		t.Fatal("expected error for empty fields")
	}
}

func TestNewHashFilter_EmptyFieldName(t *testing.T) {
	_, err := NewHashFilter([]string{"service", ""})
	if err == nil {
		t.Fatal("expected error for blank field name")
	}
}

func TestNewHashFilter_Valid(t *testing.T) {
	f, err := NewHashFilter([]string{"service", "host"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Fields()) != 2 {
		t.Errorf("expected 2 fields, got %d", len(f.Fields()))
	}
}

func TestHashFilter_NilLine(t *testing.T) {
	f, _ := NewHashFilter([]string{"svc"})
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestHashFilter_DeduplicatesSingleField(t *testing.T) {
	f, _ := NewHashFilter([]string{"service"})

	line1 := makeHashLogLine(map[string]string{"service": "auth"})
	line2 := makeHashLogLine(map[string]string{"service": "auth"})
	line3 := makeHashLogLine(map[string]string{"service": "payments"})

	if !f.Match(line1) {
		t.Error("first occurrence should match")
	}
	if f.Match(line2) {
		t.Error("duplicate should not match")
	}
	if !f.Match(line3) {
		t.Error("different value should match")
	}
}

func TestHashFilter_DeduplicatesMultipleFields(t *testing.T) {
	f, _ := NewHashFilter([]string{"service", "host"})

	a := makeHashLogLine(map[string]string{"service": "auth", "host": "h1"})
	b := makeHashLogLine(map[string]string{"service": "auth", "host": "h2"})
	c := makeHashLogLine(map[string]string{"service": "auth", "host": "h1"})

	if !f.Match(a) {
		t.Error("first combo should match")
	}
	if !f.Match(b) {
		t.Error("different host should match")
	}
	if f.Match(c) {
		t.Error("repeated combo should not match")
	}
}
