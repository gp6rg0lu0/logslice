// Package filter provides composable log line filters for logslice.
//
// BetweenFilter
//
// BetweenFilter selects log lines where a named numeric field falls within
// a closed interval [Min, Max].
//
// Usage via CLI flag --between field:min:max:
//
//	logslice --between latency:10:500 app.log
//
// Multiple --between flags may be combined; all must match (AND semantics
// when used inside a Chain).
//
// The field value is parsed as a 64-bit float, so both integer and decimal
// values are supported. Lines where the field is absent or non-numeric are
// skipped.
package filter
