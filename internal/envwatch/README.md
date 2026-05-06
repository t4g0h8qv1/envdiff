# envwatch

The `envwatch` package provides file-watching functionality for `.env` files.
It polls files at a configurable interval and emits diff events whenever a
change is detected.

## Usage

```go
import (
    "fmt"
    "time"
    "github.com/user/envdiff/internal/envwatch"
)

w := envwatch.New([]string{".env", ".env.local"}, 2*time.Second)
w.Start()
defer w.Stop()

for ev := range w.Events {
    fmt.Printf("[%s] %d change(s) in %s\n",
        ev.At.Format(time.RFC3339), len(ev.Results), ev.File)
    for _, r := range ev.Results {
        fmt.Printf("  %s  %s\n", r.Type, r.Key)
    }
}
```

## Event fields

| Field     | Type              | Description                          |
|-----------|-------------------|--------------------------------------|
| `File`    | `string`          | Path of the file that changed        |
| `At`      | `time.Time`       | Timestamp of detection               |
| `Results` | `[]diff.Result`   | Diff results between old and new     |

## Notes

- Uses polling rather than OS filesystem events for portability.
- The poll interval is configurable; a value of 1–5 seconds is recommended.
- If a watched file is deleted, a warning is printed and polling continues.
