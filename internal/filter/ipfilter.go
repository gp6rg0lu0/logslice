package filter

import (
	"fmt"
	"net"

	"github.com/nicholasgasior/logslice/internal/parser"
)

// IPFilter matches log lines where a field value falls within a given CIDR range.
type IPFilter struct {
	field string
	cidr  string
	net   *net.IPNet
}

// NewIPFilter creates a filter that matches when the given field contains an IP
// address within the specified CIDR block (e.g. "192.168.0.0/16").
func NewIPFilter(field, cidr string) (*IPFilter, error) {
	if field == "" {
		return nil, fmt.Errorf("ipfilter: field must not be empty")
	}
	if cidr == "" {
		return nil, fmt.Errorf("ipfilter: cidr must not be empty")
	}
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("ipfilter: invalid CIDR %q: %w", cidr, err)
	}
	return &IPFilter{field: field, cidr: cidr, net: network}, nil
}

// Field returns the field name used by the filter.
func (f *IPFilter) Field() string { return f.field }

// CIDR returns the CIDR string used by the filter.
func (f *IPFilter) CIDR() string { return f.cidr }

// Match returns true when the field value is a valid IP within the CIDR range.
func (f *IPFilter) Match(line *parser.LogLine) bool {
	if line == nil {
		return false
	}
	val, ok := line.Get(f.field)
	if !ok {
		return false
	}
	ip := net.ParseIP(val)
	if ip == nil {
		return false
	}
	return f.net.Contains(ip)
}
