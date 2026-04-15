package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempLog(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempLog: %v", err)
	}
	return path
}

func TestParseFlags_Defaults(t *testing.T) {
	cfg, err := parseFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != "json" {
		t.Errorf("expected default format json, got %q", cfg.Format)
	}
	if cfg.InputFile != "" {
		t.Errorf("expected empty InputFile, got %q", cfg.InputFile)
	}
}

func TestParseFlags_AllFlags(t *testing.T) {
	args := []string{
		"-from", "2024-01-01T00:00:00Z",
		"-to", "2024-01-02T00:00:00Z",
		"-level", "warn",
		"-field", "service=api",
		"-format", "text",
		"myfile.log",
	}
	cfg, err := parseFlags(args)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.From != "2024-01-01T00:00:00Z" {
		t.Errorf("From mismatch: %q", cfg.From)
	}
	if cfg.Level != "warn" {
		t.Errorf("Level mismatch: %q", cfg.Level)
	}
	if cfg.Field != "service=api" {
		t.Errorf("Field mismatch: %q", cfg.Field)
	}
	if cfg.Format != "text" {
		t.Errorf("Format mismatch: %q", cfg.Format)
	}
	if cfg.InputFile != "myfile.log" {
		t.Errorf("InputFile mismatch: %q", cfg.InputFile)
	}
}

func TestRun_InvalidLevelFilter(t *testing.T) {
	cfg := &Config{Level: "verbose"}
	if err := run(cfg); err == nil {
		t.Error("expected error for invalid level, got nil")
	}
}

func TestRun_InvalidFieldFilter(t *testing.T) {
	cfg := &Config{Field: "noequals"}
	if err := run(cfg); err == nil {
		t.Error("expected error for invalid field filter, got nil")
	}
}

func TestRun_FileNotFound(t *testing.T) {
	cfg := &Config{InputFile: "/nonexistent/path/file.log"}
	if err := run(cfg); err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestRun_BasicPipeline(t *testing.T) {
	logData := `{"level":"info","msg":"started","service":"api"}
{"level":"error","msg":"failed","service":"api"}
{"level":"debug","msg":"verbose","service":"api"}
`
	path := writeTempLog(t, logData)
	cfg := &Config{
		Level:     "error",
		InputFile: path,
		Format:    "json",
	}
	if err := run(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
