# snapshot

The `snapshot` package provides functionality to **save and compare diff results over time**.

This allows you to track how `.env` discrepancies evolve between runs, making it easier to detect regressions or verify that issues have been resolved.

## Usage

### Save a snapshot

```go
results := diff.Compare(left, right)
err := snapshot.Save("./snapshots/baseline.json", "production-baseline", results)
```

### Load a snapshot

```go
s, err := snapshot.Load("./snapshots/baseline.json")
if err != nil {
    log.Fatal(err)
}
fmt.Println("Snapshot label:", s.Label)
fmt.Println("Taken at:", s.CreatedAt)
```

### Compare against a baseline

```go
baseline, _ := snapshot.Load("./snapshots/baseline.json")
currentResults := diff.Compare(left, right)

newOrChanged := snapshot.Compare(baseline, currentResults)
if len(newOrChanged) > 0 {
    fmt.Println("New or changed issues since baseline:")
    report.Render(os.Stdout, newOrChanged, "text")
}
```

### Diff two snapshots

```go
old, _ := snapshot.Load("./snapshots/week-ago.json")
new, _ := snapshot.Load("./snapshots/today.json")

added, removed := snapshot.Diff(old, new)
fmt.Printf("%d issues resolved, %d new issues introduced\n", len(removed), len(added))
```

## Snapshot file format

Snapshots are stored as JSON:

```json
{
  "created_at": "2024-01-15T10:30:00Z",
  "label": "production-baseline",
  "results": [
    {
      "key": "DB_HOST",
      "type": "missing",
      "left_value": "localhost",
      "right_value": ""
    }
  ]
}
```
