package filter

import (
	"testing"

	"github.com/nicholasgasior/logslice/internal/parser"
)

func makeIPLogLine(field, value string) *parser.LogLine {
	fields := map[string]interface{}{}
	if field != "" {
		fields[field] = value
	}
	return parser.NewLogLine(fields)
}

func TestNewIPFilter_EmptyField(t *testing.T) {
	_, err := NewIPFilter("", "10.0.0.0/8")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewIPFilter_EmptyCIDR(t *testing.T) {
	_, err := NewIPFilter("ip", "")
	if err == nil {
		t.Fatal("expected error for empty cidr")
	}
}

func TestNewIPFilter_InvalidCIDR(t *testing.T) {
	_, err := NewIPFilter("ip", "not-a-cidr")
	if err == nil {
		t.Fatal("expected error for invalid cidr")
	}
}

func TestIPFilter_Accessors(t *testing.T) {
	f, err := NewIPFilter("src_ip", "192.168.0.0/16")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "src_ip" {
		t.Errorf("expected src_ip, got %s", f.Field())
	}
	if f.CIDR() != "192.168.0.0/16" {
		t.Errorf("expected 192.168.0.0/16, got %s", f.CIDR())
	}
}

func TestIPFilter_NilLine(t *testing.T) {
	f, _ := NewIPFilter("ip", "10.0.0.0/8")
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestIPFilter_MatchesInsideCIDR(t *testing.T) {
	f, _ := NewIPFilter("ip", "10.0.0.0/8")
	line := makeIPLogLine("ip", "10.1.2.3")
	if !f.Match(line) {
		t.Error("expected match for IP inside CIDR")
	}
}

func TestIPFilter_NoMatchOutsideCIDR(t *testing.T) {
	f, _ := NewIPFilter("ip", "10.0.0.0/8")
	line := makeIPLogLine("ip", "192.168.1.1")
	if f.Match(line) {
		t.Error("expected no match for IP outside CIDR")
	}
}

func TestIPFilter_InvalidIPValue(t *testing.T) {
	f, _ := NewIPFilter("ip", "10.0.0.0/8")
	line := makeIPLogLine("ip", "not-an-ip")
	if f.Match(line) {
		t.Error("expected no match for non-IP value")
	}
}

func TestIPFilter_MissingField(t *testing.T) {
	f, _ := NewIPFilter("ip", "10.0.0.0/8")
	line := makeIPLogLine("", "")
	if f.Match(line) {
		t.Error("expected no match when field is missing")
	}
}
