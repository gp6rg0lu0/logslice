package filter

import (
	"encoding/json"
	"strings"
)

// LogLevel represents a structured log severity level.
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelUnknown
)

// levelNames maps string representations to LogLevel values.
var levelNames = map[string]LogLevel{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"warning": LevelWarn,
	"error": LevelError,
	"err":   LevelError,
	"fatal": LevelFatal,
	"panic": LevelFatal,
}

// ParseLevel converts a string to a LogLevel, case-insensitive.
func ParseLevel(s string) LogLevel {
	if l, ok := levelNames[strings.ToLower(strings.TrimSpace(s))]; ok {
		return l
	}
	return LevelUnknown
}

// LevelFilter holds the minimum log level for filtering.
type LevelFilter struct {
	MinLevel LogLevel
}

// NewLevelFilter creates a LevelFilter from a level string.
func NewLevelFilter(minLevel string) *LevelFilter {
	return &LevelFilter{MinLevel: ParseLevel(minLevel)}
}

// Match returns true if the log line's level meets or exceeds the minimum.
// It expects a JSON-encoded log line with a "level" field.
func (f *LevelFilter) Match(line string) bool {
	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		return false
	}

	for _, key := range []string{"level", "lvl", "severity"} {
		if raw, ok := entry[key]; ok {
			if levelStr, ok := raw.(string); ok {
				parsed := ParseLevel(levelStr)
				if parsed == LevelUnknown {
					return false
				}
				return parsed >= f.MinLevel
			}
		}
	}
	return false
}
