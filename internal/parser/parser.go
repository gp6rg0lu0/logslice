// Package parser provides interfaces and implementations for parsing
// structured log lines from various input formats.
package parser

import "io"

// Parser is the interface that wraps the basic Parse method.
//
// Parse reads from r and sends parsed LogLines to the returned channel.
// The channel is closed when the reader is exhausted or an unrecoverable
// error occurs. Errors encountered during parsing are logged but do not
// stop processing; only fatal I/O errors cause the channel to close early.
type Parser interface {
	Parse(r io.Reader) <-chan *LogLine
}

// Format represents a supported log input format.
type Format string

const (
	// FormatJSON indicates newline-delimited JSON log entries.
	FormatJSON Format = "json"

	// FormatText indicates plain-text log entries matched by a regex pattern.
	FormatText Format = "text"
)

// New returns a Parser for the given format. For FormatText, pattern must
// be a named-group regular expression. For FormatJSON, pattern is ignored.
// An error is returned if the format is unknown or the pattern is invalid.
func New(format Format, pattern string) (Parser, error) {
	switch format {
	case FormatJSON:
		return NewJSONParser(), nil
	case FormatText:
		return NewTextParser(pattern)
	default:
		return nil, &UnknownFormatError{Format: string(format)}
	}
}

// UnknownFormatError is returned when an unsupported format is requested.
type UnknownFormatError struct {
	Format string
}

func (e *UnknownFormatError) Error() string {
	return "unknown log format: " + e.Format
}
