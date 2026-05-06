# audit

The `audit` package records diff run history as an append-only JSON Lines log.
Each entry captures the files compared, the full diff results, and a summary
of missing keys and conflicts.

## Usage

```go
import (
    "github.com/user/envdiff/internal/audit"
    "github.com/user/envdiff/internal/diff"
)

// After running a diff:
results := diff.Compare(left, right)
entry := audit.NewEntry([]string{".env.dev", ".env.prod"}, results)

if err := audit.Append("envdiff-audit.log", entry); err != nil {
    log.Fatal(err)
}

// Read history:
entries, err := audit.ReadAll("envdiff-audit.log")
for _, e := range entries {
    fmt.Printf("%s — missing: %d, conflicts: %d\n",
        e.Timestamp.Format(time.RFC3339),
        e.Summary.Missing,
        e.Summary.Conflicts,
    )
}
```

## Log Format

Each line in the log file is a JSON object:

```json
{"timestamp":"2024-01-15T10:30:00Z","files":[".env.dev",".env.prod"],"results":[...],"summary":{"total":5,"missing":2,"conflicts":1}}
```

## Functions

| Function | Description |
|---|---|
| `NewEntry(files, results)` | Creates a new audit entry with current timestamp |
| `Append(path, entry)` | Appends an entry to the log file (creates if absent) |
| `ReadAll(path)` | Reads and returns all entries from a log file |
