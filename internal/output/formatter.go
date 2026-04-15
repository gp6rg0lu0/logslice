package output

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Formatter converts a LogLine to a string representation.
type Formatter interface {
	Format(line *parser.LogLine) (string, error)
}

// JSONFormatter renders a LogLine as compact JSON.
type JSONFormatter struct{}

// NewJSONFormatter returns a new JSONFormatter.
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format serialises the LogLine fields to a JSON string.
func (f *JSONFormatter) Format(line *parser.LogLine) (string, error) {
	if line == nil {
		return "", fmt.Errorf("formatter: nil log line")
	}
	b, err := json.Marshal(line.Fields())
	if err != nil {
		return "", fmt.Errorf("formatter: json marshal: %w", err)
	}
	return string(b), nil
}

// TextFormatter renders a LogLine as a key=value pair string.
type TextFormatter struct {
	// keys controls output order; remaining keys follow alphabetically.
	keys []string
}

// NewTextFormatter returns a TextFormatter that prints the given keys first.
func NewTextFormatter(priorityKeys []string) *TextFormatter {
	return &TextFormatter{keys: priorityKeys}
}

// Format renders the LogLine as space-separated key=value pairs.
func (f *TextFormatter) Format(line *parser.LogLine) (string, error) {
	if line == nil {
		return "", fmt.Errorf("formatter: nil log line")
	}
	fields := line.Fields()
	seen := make(map[string]bool, len(f.keys))
	var parts []string

	for _, k := range f.keys {
		if v, ok := fields[k]; ok {
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
			seen[k] = true
		}
	}

	remaining := make([]string, 0, len(fields))
	for k := range fields {
		if !seen[k] {
			remaining = append(remaining, k)
		}
	}
	sort.Strings(remaining)
	for _, k := range remaining {
		parts = append(parts, fmt.Sprintf("%s=%v", k, fields[k]))
	}

	return strings.Join(parts, " "), nil
}
