package parser

// LogLine represents a single parsed log entry.
// Fields is a flat map of all key-value pairs extracted from the raw log line.
// The parser is responsible for populating well-known keys such as "level"
// and "time" so that downstream filters can operate on them uniformly.
type LogLine struct {
	// Fields holds every key-value pair extracted from the log entry.
	// Values are always stored as strings; numeric/boolean fields are
	// converted to their string representation during parsing.
	Fields map[string]string

	// Raw is the original, unmodified log line.
	Raw string
}

// Get returns the value associated with key, or an empty string when the
// key is absent. It is safe to call on a nil or zero-value LogLine.
func (l *LogLine) Get(key string) string {
	if l == nil || l.Fields == nil {
		return ""
	}
	return l.Fields[key]
}

// Level is a convenience accessor for the "level" field.
func (l *LogLine) Level() string {
	return l.Get("level")
}

// Time is a convenience accessor for the "time" field.
// The returned string is in whatever format was present in the log line;
// callers that need a time.Time should parse it themselves.
func (l *LogLine) Time() string {
	return l.Get("time")
}

// Parser is the interface that all log-format parsers must satisfy.
// Parse reads lines from src and sends each successfully parsed LogLine
// on the returned channel. The channel is closed when src is exhausted
// or an unrecoverable error occurs. Malformed lines are skipped.
type Parser interface {
	Parse(src interface{ Read([]byte) (int, error) }) (<-chan *LogLine, error)
}
