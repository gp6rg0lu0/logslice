package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// TextParser parses unstructured log lines using a regular expression
// with named capture groups. Each named group becomes a field in the LogLine.
type TextParser struct {
	re     *regexp.Regexp
	scanner *bufio.Scanner
	skipped int
}

// NewTextParser creates a TextParser that reads from r and parses each line
// using pattern. The pattern must be a valid Go regexp with named groups.
// Example pattern: `(?P<time>\S+) (?P<level>\w+) (?P<message>.*)`
func NewTextParser(r io.Reader, pattern string) (*TextParser, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("text parser: invalid pattern: %w", err)
	}
	names := re.SubexpNames()
	hasNamed := false
	for _, n := range names {
		if n != "" {
			hasNamed = true
			break
		}
	}
	if !hasNamed {
		return nil, fmt.Errorf("text parser: pattern must contain at least one named capture group")
	}
	return &TextParser{
		re:      re,
		scanner: bufio.NewScanner(r),
	}, nil
}

// Next returns the next parsed LogLine, or (nil, io.EOF) when the input is
// exhausted. Lines that do not match the pattern are skipped.
func (p *TextParser) Next() (*LogLine, error) {
	for p.scanner.Scan() {
		line := strings.TrimSpace(p.scanner.Text())
		if line == "" {
			continue
		}
		match := p.re.FindStringSubmatch(line)
		if match == nil {
			p.skipped++
			continue
		}
		fields := make(map[string]string)
		for i, name := range p.re.SubexpNames() {
			if name != "" && i < len(match) {
				fields[name] = match[i]
			}
		}
		return NewLogLine(fields), nil
	}
	if err := p.scanner.Err(); err != nil {
		return nil, err
	}
	return nil, io.EOF
}

// Skipped returns the number of non-empty lines that did not match the pattern
// and were skipped during parsing.
func (p *TextParser) Skipped() int {
	return p.skipped
}
