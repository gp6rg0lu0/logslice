package cli

// FormatFlag represents the output format option.
type FormatFlag string

const (
	// FormatJSON emits each matching log line as a JSON object.
	FormatJSON FormatFlag = "json"
	// FormatText emits each matching log line as plain key=value text.
	FormatText FormatFlag = "text"
)

// Config holds all parsed CLI options for a single logslice invocation.
type Config struct {
	// Input is the path to the log file, or "-" for stdin.
	Input string

	// OutputFormat controls how matched lines are written.
	OutputFormat FormatFlag

	// InputPattern is an optional regexp (with named groups) used when parsing
	// plain-text log files. When empty the JSON parser is used.
	InputPattern string

	// Level is the minimum log level to include (e.g. "warn").
	Level string

	// Since is the lower bound of the time range (RFC3339 or similar).
	Since string

	// Until is the upper bound of the time range.
	Until string

	// Fields is a slice of "key=pattern" expressions for field filtering.
	Fields []string
}

// Validate returns an error string if the Config contains obviously invalid
// combinations, or an empty string when the config looks valid.
func (c *Config) Validate() string {
	switch c.OutputFormat {
	case FormatJSON, FormatText:
		// ok
	case "":
		// default will be applied by Run
	default:
		return "unknown output format: " + string(c.OutputFormat)
	}
	return ""
}
