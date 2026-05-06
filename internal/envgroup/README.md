# envgroup

The `envgroup` package groups environment variables by prefix or custom category rules, making it easier to reason about large `.env` files.

## Functions

### `ByPrefix(env map[string]string) []Group`

Splits an env map into groups based on the first `_`-delimited segment of each key.

```
DB_HOST=localhost  →  group "DB"
DB_PORT=5432       →  group "DB"
APP_NAME=envdiff   →  group "APP"
NOPREFIX=value     →  group "_OTHER"
```

### `ByCategories(env map[string]string, categories map[string][]string) []Group`

Groups keys according to caller-supplied category definitions. Each category maps a label to one or more key prefixes.

```go
cats := map[string][]string{
    "database": {"DB_", "REDIS_"},
    "app":      {"APP_"},
}
groups := envgroup.ByCategories(env, cats)
```

Keys not matched by any category are placed in `_OTHER`.

### `Summary(groups []Group) map[string]int`

Returns a map of group name → key count, useful for quick overviews.

## Types

```go
type Group struct {
    Name string
    Keys map[string]string
}
```

## Notes

- Returned slices are always sorted alphabetically by group name.
- `_OTHER` is used as the fallback group name to avoid collisions with real prefixes.
