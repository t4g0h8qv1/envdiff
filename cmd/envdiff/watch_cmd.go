package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/envdiff/internal/envwatch"
	"github.com/user/envdiff/internal/report"
)

func runWatch(args []string) {
	fs := flag.NewFlagSet("watch", flag.ExitOnError)
	interval := fs.Duration("interval", 2*time.Second, "poll interval (e.g. 1s, 500ms)")
	format := fs.String("format", "text", "output format: text or json")
	fs.Parse(args) //nolint:errcheck

	files := fs.Args()
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "watch: at least one .env file required")
		os.Exit(1)
	}

	w := envwatch.New(files, *interval)
	w.Start()
	defer w.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Watching %d file(s) every %s. Press Ctrl+C to stop.\n",
		len(files), *interval)

	for {
		select {
		case ev := <-w.Events:
			fmt.Printf("\n[%s] Changes detected in: %s\n",
				ev.At.Format("15:04:05"), ev.File)
			out, err := report.Render(ev.Results, *format)
			if err != nil {
				fmt.Fprintf(os.Stderr, "render error: %v\n", err)
				continue
			}
			fmt.Println(out)
		case <-sig:
			fmt.Println("\nStopping watcher.")
			return
		}
	}
}
