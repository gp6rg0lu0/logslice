// Package filter provides composable log line filters for logslice.
//
// # IPFilter
//
// IPFilter matches log lines where a named field contains an IP address that
// falls within a given CIDR block.
//
// Example usage:
//
//	f, err := filter.NewIPFilter("client_ip", "10.0.0.0/8")
//	if err != nil {
//		log.Fatal(err)
//	}
//	// f.Match(line) returns true when line["client_ip"] is in 10.0.0.0/8
//
// CLI flag format:
//
//	--ip client_ip=10.0.0.0/8
//
// Multiple --ip flags are ANDed together via the filter chain.
package filter
