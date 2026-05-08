# envclone

The `envclone` package copies an env variable map into a new `.env` file, with optional filtering and redaction.

## Features

- Clone all keys or filter by **prefix** or **substring**
- **Redact** sensitive values (writes `KEY=` with empty value)
- **Overwrite protection** — refuses to clobber an existing file unless explicitly allowed
- Output is written in **sorted key order** for reproducibility

## Usage

```go
import "github.com/user/envdiff/internal/envclone"

result, err := envclone.Clone(srcMap, ".env.staging", envclone.Options{
    KeyPrefix:         "APP_",
    Redact:            false,
    OverwriteExisting: false,
})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Cloned %d keys, skipped %d\n", result.Cloned, result.Skipped)
```

## Options

| Field | Type | Description |
|---|---|---|
| `KeyPrefix` | `string` | Only clone keys with this prefix |
| `KeyContains` | `string` | Only clone keys containing this substring |
| `Redact` | `bool` | Write empty values instead of real ones |
| `OverwriteExisting` | `bool` | Allow overwriting an existing target file |

## Result

| Field | Description |
|---|---|
| `Cloned` | Number of keys written |
| `Skipped` | Number of keys excluded by filters |
| `Target` | Path of the written file |
