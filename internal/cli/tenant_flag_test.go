package cli

import (
	"testing"
)

func TestParseTenantFlag_Empty(t *testing.T) {
	f, err := ParseTenantFlag("")
	if err == nil {
		t.Fatal("expected error for empty flag, got nil")
	}
	if f != nil {
		t.Fatalf("expected nil filter, got %v", f)
	}
}

func TestParseTenantFlag_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "single tenant",
			input: "tenant_id:acme",
		},
		{
			name:  "multiple tenants",
			input: "tenant_id:acme,globex,initech",
		},
		{
			name:  "custom field",
			input: "org:alpha,beta",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := ParseTenantFlag(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f == nil {
				t.Fatal("expected non-nil filter")
			}
		})
	}
}

func TestParseTenantFlag_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "missing colon",
			input: "tenant_idacme",
		},
		{
			name:  "empty field",
			input: ":acme",
		},
		{
			name:  "empty tenants",
			input: "tenant_id:",
		},
		{
			name:  "blank tenant in list",
			input: "tenant_id:acme,,globex",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f, err := ParseTenantFlag(tc.input)
			if err == nil {
				t.Fatalf("expected error for input %q, got nil", tc.input)
			}
			if f != nil {
				t.Fatalf("expected nil filter on error, got %v", f)
			}
		})
	}
}
