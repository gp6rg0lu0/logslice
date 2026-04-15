// Package output provides types for writing and formatting parsed log lines.
//
// # Writer
//
// Writer wraps an io.Writer and serialises [parser.LogLine] values using a
// configurable [Formatter].  Use [NewWriter] to obtain a Writer that defaults
// to JSON output.
//
// # Formatters
//
// Two built-in formatters are provided:
//
//   - [JSONFormatter] – emits each log line as a compact JSON object.
//   - [TextFormatter] – emits each log line as space-separated key=value pairs.
//     A list of priority keys controls the order of the first fields; any
//     remaining fields are appended in alphabetical order.
//
// Example usage:
//
//	f := output.NewTextFormatter([]string{"time", "level", "msg"})
//	w := output.NewWriter(os.Stdout, f)
//	for _, line := range lines {
//		_ = w.Write(line)
//	}
package output
