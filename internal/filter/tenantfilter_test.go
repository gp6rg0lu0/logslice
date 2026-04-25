package filter

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeTenantLogLine(field, value string) *parser.LogLine {
	fields := map[string]interface{}{field: value}
	return parser.NewLogLine(fields)
}

func TestNewTenantFilter_EmptyField(t *testing.T) {
	_, err := NewTenantFilter("", []string{"acme"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewTenantFilter_NoTenants(t *testing.T) {
	_, err := NewTenantFilter("tenant", []string{})
	if err == nil {
		t.Fatal("expected error for empty tenant list")
	}
}

func TestNewTenantFilter_BlankTenant(t *testing.T) {
	_, err := NewTenantFilter("tenant", []string{"acme", "  "})
	if err == nil {
		t.Fatal("expected error for blank tenant ID")
	}
}

func TestTenantFilter_Accessors(t *testing.T) {
	f, err := NewTenantFilter("tenant_id", []string{"Acme", "globex"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field() != "tenant_id" {
		t.Errorf("Field() = %q, want %q", f.Field(), "tenant_id")
	}
	tenants := f.Tenants()
	if len(tenants) != 2 {
		t.Errorf("expected 2 tenants, got %d", len(tenants))
	}
}

func TestTenantFilter_NilLine(t *testing.T) {
	f, _ := NewTenantFilter("tenant", []string{"acme"})
	if f.Match(nil) {
		t.Error("expected false for nil line")
	}
}

func TestTenantFilter_MissingField(t *testing.T) {
	f, _ := NewTenantFilter("tenant", []string{"acme"})
	line := makeTenantLogLine("other", "acme")
	if f.Match(line) {
		t.Error("expected false when field is missing")
	}
}

func TestTenantFilter_MatchExact(t *testing.T) {
	f, _ := NewTenantFilter("tenant", []string{"acme", "globex"})
	line := makeTenantLogLine("tenant", "acme")
	if !f.Match(line) {
		t.Error("expected match for tenant 'acme'")
	}
}

func TestTenantFilter_MatchCaseInsensitive(t *testing.T) {
	f, _ := NewTenantFilter("tenant", []string{"Acme"})
	line := makeTenantLogLine("tenant", "ACME")
	if !f.Match(line) {
		t.Error("expected case-insensitive match")
	}
}

func TestTenantFilter_NoMatch(t *testing.T) {
	f, _ := NewTenantFilter("tenant", []string{"acme"})
	line := makeTenantLogLine("tenant", "initech")
	if f.Match(line) {
		t.Error("expected no match for unknown tenant")
	}
}
