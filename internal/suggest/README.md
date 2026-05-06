# suggest

The `suggest` package analyzes diff results and generates human-readable recommendations for resolving missing keys and conflicting values across environment files.

## Features

- Generates actionable suggestions for `missing` and `conflict` diff results
- Provides formatted output suitable for CLI display

## Usage

```go
import (
    "fmt"
    "github.com/user/envdiff/internal/diff"
    "github.com/user/envdiff/internal/suggest"
)

results := diff.Compare(envMaps)
suggestions := suggest.Generate(results)
fmt.Print(suggest.Format(suggestions))
```

## Output Example

```
Suggestions:
  1. [MISSING] Add key "DB_HOST" to .env.production
  2. [CONFLICT] Resolve conflicting values for key "API_URL" across: .env.staging, .env.production
```

If no issues are found:

```
No suggestions — all environments look consistent.
```
