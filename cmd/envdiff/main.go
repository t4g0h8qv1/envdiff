package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/filter"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/report"
)

func main() {
	format := flag.String("format", "text", "Output format: text or json")
	onlyMissing := flag.Bool("missing", false, "Show only missing keys")
	onlyConflicts := flag.Bool("conflicts", false, "Show only conflicting keys")
	keyPrefix := flag.String("prefix", "", "Filter keys by prefix")
	keyContains := flag.String("contains", "", "Filter keys containing substring")
	flag.Parse()

	paths := flag.Args()
	if len(paths) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: envdiff [flags] <file1> <file2>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	envMaps, err := loader.LoadFiles(paths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading files: %v\n", err)
		os.Exit(1)
	}

	results := diff.Compare(envMaps[0], envMaps[1])

	filterOpts := filter.Options{
		OnlyMissing:   *onlyMissing,
		OnlyConflicts: *onlyConflicts,
		KeyPrefix:     *keyPrefix,
		KeyContains:   *keyContains,
	}
	results = filter.Apply(results, filterOpts)

	output, err := report.Render(results, *format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering report: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(output)
}
