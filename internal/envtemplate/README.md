# envtemplate

The `envtemplate` package generates a `.env.template` (or `.env.example`) file
from one or more parsed environment maps.

All **keys** are preserved and sorted alphabetically. **Values** are replaced
with an optional placeholder string so that the template can be safely committed
to version control without leaking secrets.

## Usage

```go
import "github.com/yourorg/envdiff/internal/envtemplate"

opts := envtemplate.Options{
    Placeholder:     "CHANGE_ME",
    IncludeComments: true,
}

// Generate returns []string of KEY=PLACEHOLDER lines.
lines := envtemplate.Generate(envMaps, opts)

// Write saves the template directly to a file.
err := envtemplate.Write(".env.template", envMaps, opts)
```

## Options

| Field | Type | Description |
|---|---|---|
| `Placeholder` | `string` | Value written for every key. Defaults to empty (`KEY=`). |
| `IncludeComments` | `bool` | Prepends a header comment warning against storing real secrets. |

## Behaviour

- Keys from **all** supplied env maps are merged into a single deduplicated set.
- Output lines are **sorted alphabetically** by key for stable diffs.
- A trailing newline is always appended when writing to a file.
