// Package filter provides composable log-line filters used by logslice.
//
// PathFilter
//
// PathFilter checks whether the string value of a named field matches one or
// more URL path patterns.  Two matching modes are supported:
//
//   - Prefix mode (default): the field value must start with at least one of
//     the configured paths.  Useful for grouping an entire sub-tree of routes,
//     e.g. "/api/v1" matches "/api/v1/users" and "/api/v1/orders".
//
//   - Exact mode: the field value must equal one of the configured paths
//     exactly.  Useful for health-check or readiness endpoints where only the
//     literal path should be selected.
//
// Construction
//
//	f, err := filter.NewPathFilter("path", []string{"/api/v1", "/health"}, false)
//
// CLI flag format
//
//	--path url:/api/v1,/health          # prefix mode
//	--path url:exact:/ready,/healthz    # exact mode
package filter
