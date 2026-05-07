# envpromote

The `envpromote` package copies environment variables from a **source** map (e.g. staging) into a **destination** map (e.g. production) according to configurable rules.

## Features

- **AllowKeys** – promote only an explicit set of keys
- **DenyKeys** – exclude specific keys from promotion (e.g. secrets)
- **Conflict policies** – choose how to handle keys that already exist in the destination:
  - `skip` (default) – leave the destination value unchanged
  - `overwrite` – replace the destination value with the source value
  - `error` – abort and return an error

## Usage

```go
import "github.com/yourorg/envdiff/internal/envpromote"

src := map[string]string{"DB_HOST": "staging-db", "SECRET": "abc"}
dst := map[string]string{"DB_HOST": "prod-db"}

results, err := envpromote.Promote(src, dst, envpromote.Options{
    DenyKeys: []string{"SECRET"},
    Policy:   envpromote.PolicySkip,
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(envpromote.Summary(results))
// Output: 0 promoted, 1 skipped
```

## Result Actions

| Action        | Meaning                                          |
|---------------|--------------------------------------------------|
| `promoted`    | Key did not exist in dst; value was copied       |
| `overwritten` | Key existed in dst; value was replaced (policy=overwrite) |
| `skipped`     | Key existed in dst; value was left unchanged     |

## Notes

- `Promote` mutates the `dst` map in place.
- Use `AllowKeys` and `DenyKeys` together to fine-tune which variables cross environment boundaries.
