# ignore

The `ignore` package provides support for `.envignore` files — a simple way to
exclude specific keys or key prefixes from `envdiff` comparisons.

## File Format

Create a `.envignore` file in your project root (or pass it via CLI flag).
Each line specifies a key or pattern to ignore:

```
# Lines starting with '#' are comments

# Exact key match
SECRET_KEY
API_TOKEN

# Prefix wildcard — ignores all keys starting with AWS_
AWS_*
INTERNAL_*
```

## Usage

```go
rules, err := ignore.LoadFile(".envignore")
if err != nil {
    log.Fatal(err)
}

// Check a single key
if rules.ShouldIgnore("SECRET_KEY") {
    fmt.Println("skipping")
}

// Strip ignored keys from a parsed env map before diffing
clean := rules.FilterKeys(envMap)
```

## Integration

The loader and diff pipeline accept pre-filtered maps. Use `FilterKeys` after
parsing each environment file and before calling `diff.Compare` to exclude
sensitive or irrelevant keys from the report.
