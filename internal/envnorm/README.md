# envnorm

The `envnorm` package normalizes `.env` file key/value maps according to configurable conventions.

## Features

- **Uppercase keys** — converts all keys to `UPPER_SNAKE_CASE`
- **Trim values** — strips leading/trailing whitespace from values
- **Remove empty** — optionally drops keys with empty values
- **Quote values** — optionally wraps unquoted values in double quotes

## Usage

```go
import "github.com/yourorg/envdiff/internal/envnorm"

env := map[string]string{
    "db_host": "  localhost  ",
    "api_key":  "secret",
    "EMPTY":    "",
}

opts := envnorm.DefaultOptions()
opts.RemoveEmpty = true

normalized, violations := envnorm.Normalize(env, opts)

for _, v := range violations {
    fmt.Printf("[%s] %s → %s (%s)\n", v.Key, v.Original, v.Normalized, v.Reason)
}
```

## Options

| Option         | Default | Description                              |
|----------------|---------|------------------------------------------|
| `UppercaseKeys`| `true`  | Convert keys to uppercase                |
| `TrimValues`   | `true`  | Strip whitespace from values             |
| `RemoveEmpty`  | `false` | Drop keys whose value is empty           |
| `QuoteValues`  | `false` | Wrap unquoted values in double quotes    |

## Violations

Each `Violation` records the key, original value, normalized value, and a human-readable reason string, making it easy to audit what changed.
