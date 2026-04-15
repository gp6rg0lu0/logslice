// Package cli wires together flags, filters, parser, and output writer
// to implement the logslice command-line interface.
package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
)

// Config holds the parsed CLI flags.
type Config struct {
	From      string
	To        string
	Level     string
	Field     string
	Format    string
	InputFile string
}

// Run parses args and executes the main pipeline.
func Run(args []string) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}
	return run(cfg)
}

func parseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	cfg := &Config{}

	fs.StringVar(&cfg.From, "from", "", "start time (RFC3339 or epoch ms)")
	fs.StringVar(&cfg.To, "to", "", "end time (RFC3339 or epoch ms)")
	fs.StringVar(&cfg.Level, "level", "", "minimum log level (debug|info|warn|error)")
	fs.StringVar(&cfg.Field, "field", "", "field filter in key=value format")
	fs.StringVar(&cfg.Format, "format", "json", "output format: json or text")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() > 0 {
		cfg.InputFile = fs.Arg(0)
	}
	return cfg, nil
}

func run(cfg *Config) error {
	// Build filter chain.
	var filters []filter.Filter

	if cfg.From != "" || cfg.To != "" {
		tf, err := filter.NewTimeFilter(cfg.From, cfg.To)
		if err != nil {
			return fmt.Errorf("time filter: %w", err)
		}
		filters = append(filters, tf)
	}

	if cfg.Level != "" {
		lf, err := filter.NewLevelFilter(cfg.Level)
		if err != nil {
			return fmt.Errorf("level filter: %w", err)
		}
		filters = append(filters, lf)
	}

	if cfg.Field != "" {
		ff, err := filter.NewFieldFilter(cfg.Field)
		if err != nil {
			return fmt.Errorf("field filter: %w", err)
		}
		filters = append(filters, ff)
	}

	chain := filter.NewChain(filters...)

	// Open input.
	var src io.Reader = os.Stdin
	if cfg.InputFile != "" {
		f, err := os.Open(cfg.InputFile)
		if err != nil {
			return fmt.Errorf("open input: %w", err)
		}
		defer f.Close()
		src = f
	}

	p := parser.NewJSONParser(src)
	w := output.NewWriter(os.Stdout, cfg.Format)

	for {
		line, err := p.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("parse: %w", err)
		}
		if chain.Match(line) {
			if werr := w.Write(line); werr != nil {
				return fmt.Errorf("write: %w", werr)
			}
		}
	}
	return nil
}
