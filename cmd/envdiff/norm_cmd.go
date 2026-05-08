package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envdiff/internal/envnorm"
	"github.com/yourorg/envdiff/internal/parser"
)

func runNorm(args []string) {
	fs := flag.NewFlagSet("norm", flag.ExitOnError)
	uppercase := fs.Bool("uppercase", true, "Uppercase all keys")
	trim := fs.Bool("trim", true, "Trim whitespace from values")
	removeEmpty := fs.Bool("remove-empty", false, "Remove keys with empty values")
	quote := fs.Bool("quote", false, "Quote unquoted values")
	formatFlag := fs.String("format", "text", "Output format: text or json")
	dryRun := fs.Bool("dry-run", true, "Print changes without writing (default true)")
	_ = fs.Parse(args)

	paths := fs.Args()
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "usage: envdiff norm [options] <file...>")
		os.Exit(1)
	}

	opts := envnorm.Options{
		UppercaseKeys: *uppercase,
		TrimValues:    *trim,
		RemoveEmpty:   *removeEmpty,
		QuoteValues:   *quote,
	}

	type fileResult struct {
		File       string              `json:"file"`
		Violations []envnorm.Violation `json:"violations"`
	}

	var results []fileResult

	for _, path := range paths {
		env, err := parser.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
			os.Exit(1)
		}

		normalized, violations := envnorm.Normalize(env, opts)
		results = append(results, fileResult{File: path, Violations: violations})

		if !*dryRun && len(violations) > 0 {
			f, err := os.Create(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error writing %s: %v\n", path, err)
				os.Exit(1)
			}
			for k, v := range normalized {
				fmt.Fprintf(f, "%s=%s\n", k, v)
			}
			f.Close()
		}
	}

	if *formatFlag == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(results)
		return
	}

	for _, r := range results {
		if len(r.Violations) == 0 {
			fmt.Printf("%s: already normalized\n", r.File)
			continue
		}
		fmt.Printf("%s:\n", r.File)
		for _, v := range r.Violations {
			fmt.Printf("  [%s] %s → %s (%s)\n", v.Key, v.Original, v.Normalized, v.Reason)
		}
	}
}
