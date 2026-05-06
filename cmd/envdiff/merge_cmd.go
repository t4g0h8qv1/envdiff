package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/yourorg/envdiff/internal/merge"
	"github.com/yourorg/envdiff/internal/parser"
)

// runMerge handles the "merge" sub-command.
// Usage: envdiff merge -strategy=left file1.env file2.env
func runMerge(args []string) error {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	strategy := fs.String("strategy", "left", "merge strategy: left | right | union")
	output := fs.String("output", "env", "output format: env | json")

	if err := fs.Parse(args); err != nil {
		return err
	}

	paths := fs.Args()
	if len(paths) != 2 {
		return fmt.Errorf("merge requires exactly 2 env files, got %d", len(paths))
	}

	left, err := parser.ParseFile(paths[0])
	if err != nil {
		return fmt.Errorf("parsing %s: %w", paths[0], err)
	}

	right, err := parser.ParseFile(paths[1])
	if err != nil {
		return fmt.Errorf("parsing %s: %w", paths[1], err)
	}

	result, err := merge.Merge(left, right, merge.Strategy(*strategy))
	if err != nil {
		return err
	}

	for _, w := range result.Warnings {
		fmt.Fprintf(os.Stderr, "WARN: %s\n", w)
	}

	switch *output {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result.Env)
	case "env":
		keys := make([]string, 0, len(result.Env))
		for k := range result.Env {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Printf("%s=%s\n", k, result.Env[k])
		}
		return nil
	default:
		return fmt.Errorf("unknown output format: %q", *output)
	}
}
