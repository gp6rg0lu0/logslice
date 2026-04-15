// Package parser provides readers that decode structured log lines into
// field maps consumable by the filter package.
package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// Entry represents a single decoded log line as a flat string map.
type Entry map[string]string

// JSONParser reads newline-delimited JSON log files and decodes each line
// into an Entry. Non-string JSON values are converted to their string
// representation so that downstream filters always operate on strings.
type JSONParser struct {
	scanner *bufio.Scanner
}

// NewJSONParser wraps r in a JSONParser.
func NewJSONParser(r io.Reader) *JSONParser {
	return &JSONParser{scanner: bufio.NewScanner(r)}
}

// Next advances to the next log line and returns the decoded Entry.
// It returns (nil, io.EOF) when the input is exhausted.
func (p *JSONParser) Next() (Entry, error) {
	for p.scanner.Scan() {
		line := p.scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var raw map[string]interface{}
		if err := json.Unmarshal(line, &raw); err != nil {
			return nil, fmt.Errorf("parse error: %w", err)
		}
		entry := make(Entry, len(raw))
		for k, v := range raw {
			switch val := v.(type) {
			case string:
				entry[k] = val
			default:
				entry[k] = fmt.Sprintf("%v", val)
			}
		}
		return entry, nil
	}
	if err := p.scanner.Err(); err != nil {
		return nil, err
	}
	return nil, io.EOF
}
