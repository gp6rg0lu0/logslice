package cli

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/filter"
)

// ParseTenantFlag parses a --tenant flag value of the form:
//
//	"<field>:<id1>,<id2>,..."
//
// Example: "tenant_id:acme,globex,initech"
func ParseTenantFlag(value string) (*filter.TenantFilter, error) {
	if value == "" {
		return nil, nil
	}

	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("--tenant: expected format <field>:<id1>,<id2>,... got %q", value)
	}

	field := strings.TrimSpace(parts[0])
	if field == "" {
		return nil, fmt.Errorf("--tenant: field name must not be empty")
	}

	raw := strings.Split(parts[1], ",")
	tenants := make([]string, 0, len(raw))
	for _, t := range raw {
		trimmed := strings.TrimSpace(t)
		if trimmed == "" {
			continue
		}
		tenants = append(tenants, trimmed)
	}

	if len(tenants) == 0 {
		return nil, fmt.Errorf("--tenant: at least one tenant ID is required")
	}

	f, err := filter.NewTenantFilter(field, tenants)
	if err != nil {
		return nil, fmt.Errorf("--tenant: %w", err)
	}
	return f, nil
}
