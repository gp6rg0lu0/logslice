// Package output handles formatting and writing filtered log lines.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Format represents the output format for log lines.
type Format string

const (
	// FormatJSON outputs each log line as compact JSON.
	FormatJSON Format = "json"
	// FormatText outputs each log line as key=value pairs.
	FormatText Format = "text"
	// FormatRaw outputs the original raw line.
	FormatRaw Format = "raw"
)

// Writer writes log lines to an io.Writer in the specified format.
type Writer struct {
	w      io.Writer
	format Format
}

// NewWriter creates a new Writer with the given output destination and format.
// If format is unrecognized, FormatJSON is used.
func NewWriter(w io.Writer, format Format) *Writer {
	switch format {
	case FormatJSON, FormatText, FormatRaw:
		// valid
	default:
		format = FormatJSON
	}
	return &Writer{w: w, format: format}
}

// Write outputs a single log line according to the configured format.
func (wr *Writer) Write(line *parser.LogLine) error {
	if line == nil {
		return nil
	}
	switch wr.format {
	case FormatText:
		return wr.writeText(line)
	case FormatRaw:
		_, err := fmt.Fprintln(wr.w, line.Raw())
		return err
	default:
		return wr.writeJSON(line)
	}
}

func (wr *Writer) writeJSON(line *parser.LogLine) error {
	b, err := json.Marshal(line.Fields())
	if err != nil {
		return fmt.Errorf("output: marshal json: %w", err)
	}
	_, err = fmt.Fprintf(wr.w, "%s\n", b)
	return err
}

func (wr *Writer) writeText(line *parser.LogLine) error {
	parts := make([]string, 0, len(line.Fields()))
	for k, v := range line.Fields() {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	_, err := fmt.Fprintln(wr.w, strings.Join(parts, " "))
	return err
}
