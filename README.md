# logslice

A fast CLI tool for filtering and slicing structured log files by time range, level, or field pattern.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

---

## Usage

```bash
# Filter logs by time range
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" app.log

# Filter by log level
logslice --level error app.log

# Filter by field pattern
logslice --field "user_id=42" app.log

# Combine filters
logslice --from "2024-01-15T08:00:00Z" --level warn --field "service=api" app.log

# Read from stdin
cat app.log | logslice --level error
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start of time range (RFC3339) |
| `--to` | End of time range (RFC3339) |
| `--level` | Minimum log level (`debug`, `info`, `warn`, `error`) |
| `--field` | Field pattern to match (`key=value`) |
| `--format` | Input format: `json`, `logfmt` (default: `json`) |

---

## Requirements

- Go 1.21+

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)