# envscope

Provides scoped environment variable resolution across multiple named environments.

## Overview

`envscope` lets you define an ordered list of named scopes (e.g. `prod`, `staging`, `dev`) and resolve keys with a clear priority: the first scope that defines a key wins.

## Usage

```go
import "github.com/user/envdiff/internal/envscope"

prod := envscope.Scope{
    Name: "prod",
    Vars: map[string]string{"DB_HOST": "prod-db", "PORT": "443"},
}
dev := envscope.Scope{
    Name: "dev",
    Vars: map[string]string{"DB_HOST": "localhost", "DEBUG": "true"},
}

r := envscope.New(prod, dev) // prod has highest priority

// Resolve a single key
val, fromScope, err := r.Resolve("DB_HOST")
// val = "prod-db", fromScope = "prod"

// Merge all scopes (higher priority wins)
all := r.ResolveAll()

// Detect keys with differing values across scopes
conflicts := r.FindConflicts()
for _, c := range conflicts {
    fmt.Printf("conflict on %s: %v\n", c.Key, c.Scopes)
}
```

## API

| Function | Description |
|---|---|
| `New(scopes ...Scope)` | Create a resolver with ordered scopes |
| `Resolve(key)` | Return value + source scope for a key |
| `ResolveAll()` | Merge all scopes into a single map |
| `FindConflicts()` | Return keys with differing values across scopes |
