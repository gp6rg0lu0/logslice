package cli

import "testing"

func TestConfig_Validate_Defaults(t *testing.T) {
	c := &Config{}
	if msg := c.Validate(); msg != "" {
		t.Errorf("expected empty validation message, got %q", msg)
	}
}

func TestConfig_Validate_JSONFormat(t *testing.T) {
	c := &Config{OutputFormat: FormatJSON}
	if msg := c.Validate(); msg != "" {
		t.Errorf("expected empty validation message, got %q", msg)
	}
}

func TestConfig_Validate_TextFormat(t *testing.T) {
	c := &Config{OutputFormat: FormatText}
	if msg := c.Validate(); msg != "" {
		t.Errorf("expected empty validation message, got %q", msg)
	}
}

func TestConfig_Validate_UnknownFormat(t *testing.T) {
	c := &Config{OutputFormat: "xml"}
	msg := c.Validate()
	if msg == "" {
		t.Fatal("expected validation error for unknown format")
	}
}

func TestFormatFlag_Constants(t *testing.T) {
	if FormatJSON != "json" {
		t.Errorf("FormatJSON should be %q, got %q", "json", FormatJSON)
	}
	if FormatText != "text" {
		t.Errorf("FormatText should be %q, got %q", "text", FormatText)
	}
}

func TestConfig_Fields(t *testing.T) {
	c := &Config{
		Fields: []string{"level=ERROR", "service=api"},
	}
	if len(c.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(c.Fields))
	}
	if c.Fields[0] != "level=ERROR" {
		t.Errorf("unexpected field: %q", c.Fields[0])
	}
}
