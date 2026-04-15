// Package parser provides log line parsing for logslice.
//
// # Supported Formats
//
// Two formats are currently supported:
//
//   - JSON (FormatJSON): newline-delimited JSON objects. Each line is
//     unmarshalled into a map[string]string; non-string values are
//     converted to their string representation.
//
//   - Text (FormatText): plain-text lines matched against a named-group
//     regular expression supplied by the caller. Each named capture group
//     becomes a field on the resulting LogLine.
//
// # Usage
//
//	p, err := parser.New(parser.FormatJSON, "")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for line := range p.Parse(os.Stdin) {
//		fmt.Println(line.Level())
//	}
package parser
