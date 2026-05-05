package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/loader"
	"github.com/user/envdiff/internal/report"
)

func main() {
	format := flag.String("format", "text", "Output format: text or json")
	dir := flag.String("dir", "", "Directory to scan for .env files")
	flag.Parse()

	var files []loader.EnvFile
	var err error

	args := flag.Args()
	if *dir != "" {
		files, err = loader.LoadDir(*dir)
	} else if len(args) >= 2 {
		files, err = loader.LoadFiles(args)
	} else {
		fmt.Fprintln(os.Stderr, "Usage: envdiff [--format text|json] [--dir <dir>] <file1> <file2> [...]")
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading files: %v\n", err)
		os.Exit(1)
	}

	if len(files) < 2 {
		fmt.Fprintln(os.Stderr, "At least two .env files are required for comparison")
		os.Exit(1)
	}

	for i := 0; i < len(files)-1; i++ {
		left := files[i]
		right := files[i+1]

		results := diff.Compare(left.Entries, right.Entries)

		output, err := report.Render(results, left.Name, right.Name, *format)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error rendering report: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(output)
	}
}
