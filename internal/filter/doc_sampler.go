// Package filter provides log line filtering primitives.
//
// SamplerFilter
//
// SamplerFilter passes every Nth log line, useful for reducing volume
// when tailing high-throughput logs. The counter is 1-based so the
// first matching line is line number N, then 2N, 3N, etc.
//
// Example usage:
//
//	f, err := filter.NewSamplerFilter(10) // keep every 10th line
//	if err != nil { ... }
//	chain := filter.NewChain(levelFilter, f)
//
// The filter is safe for concurrent use via atomic counter.
package filter
