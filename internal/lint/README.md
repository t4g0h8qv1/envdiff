# lint

The `lint` package provides static analysis rules for `.env` files, helping teams enforce consistent key naming and value conventions.

## Rules

| Rule | Description |
|------|-------------|
| `uppercase_keys` | All keys must be fully uppercase |
| `no_spaces_in_keys` | Keys must not contain spaces |
| `no_empty_values` | Values must not be empty or whitespace-only |
| `key_format` | Keys must match `[A-Z][A-Z0-9_]*` |

## Usage

```go
import "github.com/yourorg/envdiff/internal/lint"

envMaps := map[string]map[string]string{
    ".env.staging": {
        "database_url": "postgres://localhost/db",
        "PORT":         "",
    },
}

rules := []lint.Rule{
    lint.RuleUppercaseKeys,
    lint.RuleNoEmptyValues,
    lint.RuleKeyFormat,
}

violations := lint.Check(envMaps, rules)
for _, v := range violations {
    fmt.Printf("[%s] %s — %s\n", v.File, v.Key, v.Message)
}
```

## CLI Integration

The `envdiff lint` subcommand runs all enabled rules against one or more env files and reports violations with file and key context.

## Exit Codes

- `0` — No violations found
- `1` — One or more violations detected
