# envdiff

Compare `.env` files across environments and highlight missing or conflicting keys.

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git && cd envdiff && go build ./...
```

## Usage

```bash
envdiff [flags] <file1> <file2> [file3...]
```

### Example

```bash
envdiff .env.development .env.production
```

**Output:**

```
MISSING in .env.production:
  - DATABASE_URL
  - REDIS_HOST

CONFLICT:
  - API_BASE_URL
      .env.development  → http://localhost:3000
      .env.production   → https://api.example.com
```

### Flags

| Flag | Description |
|------|-------------|
| `--missing` | Show only missing keys |
| `--conflicts` | Show only conflicting values |
| `--quiet` | Exit with non-zero status if differences found (useful in CI) |

### CI Integration

```bash
# Fail the pipeline if environments are out of sync
envdiff --quiet .env.example .env.staging
```

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

[MIT](LICENSE)