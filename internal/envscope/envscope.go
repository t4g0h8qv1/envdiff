// Package envscope provides scoped environment variable resolution,
// allowing values to be looked up across multiple named environments
// with a defined priority order.
package envscope

import "fmt"

// Scope represents a named set of environment variables.
type Scope struct {
	Name string
	Vars map[string]string
}

// Resolver holds an ordered list of scopes and resolves keys by priority.
type Resolver struct {
	scopes []Scope
}

// New creates a Resolver with the given scopes in priority order (first = highest).
func New(scopes ...Scope) *Resolver {
	return &Resolver{scopes: scopes}
}

// Resolve returns the value for key from the highest-priority scope that defines it.
// The second return value is the name of the scope the value came from.
// Returns an error if no scope defines the key.
func (r *Resolver) Resolve(key string) (string, string, error) {
	for _, s := range r.scopes {
		if val, ok := s.Vars[key]; ok {
			return val, s.Name, nil
		}
	}
	return "", "", fmt.Errorf("key %q not found in any scope", key)
}

// ResolveAll returns a map of all keys resolved across scopes.
// Each key is resolved from the highest-priority scope that defines it.
func (r *Resolver) ResolveAll() map[string]string {
	result := make(map[string]string)
	// iterate in reverse so higher-priority scopes overwrite lower ones
	for i := len(r.scopes) - 1; i >= 0; i-- {
		for k, v := range r.scopes[i].Vars {
			result[k] = v
		}
	}
	return result
}

// Conflicts returns keys that are defined in more than one scope with different values.
type Conflict struct {
	Key    string
	Scopes map[string]string // scope name -> value
}

// FindConflicts returns all keys that have differing values across scopes.
func (r *Resolver) FindConflicts() []Conflict {
	// collect all keys and their per-scope values
	keyScopes := make(map[string]map[string]string)
	for _, s := range r.scopes {
		for k, v := range s.Vars {
			if keyScopes[k] == nil {
				keyScopes[k] = make(map[string]string)
			}
			keyScopes[k][s.Name] = v
		}
	}

	var conflicts []Conflict
	for key, sv := range keyScopes {
		if len(sv) < 2 {
			continue
		}
		var first string
		diverges := false
		for _, v := range sv {
			if first == "" {
				first = v
			} else if v != first {
				diverges = true
				break
			}
		}
		if diverges {
			conflicts = append(conflicts, Conflict{Key: key, Scopes: sv})
		}
	}
	return conflicts
}

// ScopesForKey returns the names of all scopes that define the given key,
// in priority order (highest priority first).
func (r *Resolver) ScopesForKey(key string) []string {
	var names []string
	for _, s := range r.scopes {
		if _, ok := s.Vars[key]; ok {
			names = append(names, s.Name)
		}
	}
	return names
}
