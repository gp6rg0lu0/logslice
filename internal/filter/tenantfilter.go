package filter

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// TenantFilter matches log lines where a field value equals one of the
// allowed tenant identifiers. Comparison is case-insensitive by default.
type TenantFilter struct {
	field   string
	tenants map[string]struct{}
}

// NewTenantFilter creates a TenantFilter that matches lines whose field value
// is one of the provided tenant IDs. Field and each tenant must be non-empty.
func NewTenantFilter(field string, tenants []string) (*TenantFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("tenantfilter: field must not be empty")
	}
	if len(tenants) == 0 {
		return nil, fmt.Errorf("tenantfilter: at least one tenant ID is required")
	}
	set := make(map[string]struct{}, len(tenants))
	for _, t := range tenants {
		if strings.TrimSpace(t) == "" {
			return nil, fmt.Errorf("tenantfilter: tenant ID must not be blank")
		}
		set[strings.ToLower(t)] = struct{}{}
	}
	return &TenantFilter{field: field, tenants: set}, nil
}

// Field returns the log field inspected by this filter.
func (f *TenantFilter) Field() string { return f.field }

// Tenants returns the set of accepted tenant IDs (lower-cased).
func (f *TenantFilter) Tenants() []string {
	out := make([]string, 0, len(f.tenants))
	for t := range f.tenants {
		out = append(out, t)
	}
	return out
}

// Match returns true when the line's field value (case-insensitive) is one of
// the configured tenant IDs.
func (f *TenantFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	v, ok := line.Get(f.field)
	if !ok {
		return false
	}
	_, found := f.tenants[strings.ToLower(v)]
	return found
}
