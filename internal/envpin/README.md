# envpin

The `envpin` package lets you **pin** specific environment variable keys to expected values and detect when those values drift.

## Use Cases

- Lock critical keys (e.g. `APP_ENV=production`) and alert when they change.
- CI checks that prevent accidental environment promotions.
- Audit trail of when a key was pinned.

## API

### `Save(path string, env map[string]string, keys []string) error`

Pins the given keys from `env` to a JSON file at `path`. Returns an error if any key is absent from the map.

### `Load(path string) (PinFile, error)`

Reads a previously saved pin file from disk.

### `Check(pf PinFile, env map[string]string) []Violation`

Compares the current `env` map against the loaded `PinFile`. Returns a slice of `Violation` for any key whose value has changed or is missing.

## Types

```go
type Pin struct {
    Key      string
    Expected string
    PinnedAt string
}

type Violation struct {
    Key     string
    Pinned  string
    Actual  string
    Missing bool
}
```

## Example

```go
// Pin current production values
envpin.Save("pins.json", prodEnv, []string{"APP_ENV", "DB_HOST"})

// Later, verify staging matches
pf, _ := envpin.Load("pins.json")
violations := envpin.Check(pf, stagingEnv)
for _, v := range violations {
    fmt.Printf("DRIFT: %s pinned=%s actual=%s\n", v.Key, v.Pinned, v.Actual)
}
```
