# merge

The `merge` package provides functionality to combine two `.env` file maps into a single unified map using a configurable conflict-resolution strategy.

## Strategies

| Strategy | Behaviour |
|----------|-----------|
| `left`   | On conflict, the **left** file's value wins. |
| `right`  | On conflict, the **right** file's value wins. |
| `union`  | All keys from both files are included. Conflicts keep the **left** value and emit a warning. |

## Usage

```go
import "github.com/yourorg/envdiff/internal/merge"

left  := map[string]string{"DB_HOST": "localhost", "PORT": "5432"}
right := map[string]string{"DB_HOST": "prod-db",   "LOG_LEVEL": "info"}

result, err := merge.Merge(left, right, merge.StrategyUnion)
if err != nil {
    log.Fatal(err)
}

for k, v := range result.Env {
    fmt.Printf("%s=%s\n", k, v)
}

for _, w := range result.Warnings {
    fmt.Println("WARN:", w)
}
```

## Warnings

Whenever a conflict is detected, a human-readable warning is appended to `Result.Warnings`. Warnings are sorted alphabetically by key name for deterministic output.
