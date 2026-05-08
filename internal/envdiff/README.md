# envdiff

The `envdiff` package provides drift detection across two or more `.env` environments.

## Features

- **Two-environment diff** (`Detect`): compare a pair of env maps and classify each key as missing or conflicting, with severity levels based on key patterns.
- **Multi-environment drift** (`DetectMulti`): compare any number of named environments at once and produce a `DriftReport`.
- **Formatted output** (`Print`, `PrintDrift`): render results as human-readable text (with ANSI colour) or structured JSON.

## Usage

### Two environments

```go
results := envdiff.Detect(mapA, mapB)
envdiff.Print(os.Stdout, results, "text")
```

### Multiple environments

```go
envs := map[string]map[string]string{
    "dev":     loadEnv("dev.env"),
    "staging": loadEnv("staging.env"),
    "prod":    loadEnv("prod.env"),
}
report := envdiff.DetectMulti(envs)
fmt.Println(report.Summary())
envdiff.PrintDrift(os.Stdout, report, "json")
```

## DriftEntry statuses

| Status     | Meaning                                          |
|------------|--------------------------------------------------|
| `missing`  | Key present in at least one env but not all      |
| `conflict` | Key present everywhere but values differ         |

## Output formats

- `text` — aligned table with ANSI colour coding (default)
- `json` — machine-readable JSON suitable for CI pipelines
