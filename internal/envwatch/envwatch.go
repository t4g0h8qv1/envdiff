// Package envwatch watches .env files for changes and emits diff events.
package envwatch

import (
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
)

// Event represents a change detected in a watched .env file.
type Event struct {
	File    string
	At      time.Time
	Results []diff.Result
}

// Watcher polls one or more .env files for changes.
type Watcher struct {
	files    []string
	interval time.Duration
	prev     map[string]map[string]string
	Events   chan Event
	stop     chan struct{}
}

// New creates a new Watcher for the given files and poll interval.
func New(files []string, interval time.Duration) *Watcher {
	return &Watcher{
		files:    files,
		interval: interval,
		prev:     make(map[string]map[string]string),
		Events:   make(chan Event, 16),
		stop:     make(chan struct{}),
	}
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	for _, f := range w.files {
		if m, err := parser.ParseFile(f); err == nil {
			w.prev[f] = m
		}
	}
	go w.loop()
}

// Stop signals the watcher to halt.
func (w *Watcher) Stop() {
	close(w.stop)
}

func (w *Watcher) loop() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-w.stop:
			return
		case <-ticker.C:
			w.poll()
		}
	}
}

func (w *Watcher) poll() {
	for _, f := range w.files {
		curr, err := parser.ParseFile(f)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("envwatch: file removed: %s\n", f)
			}
			continue
		}
		prev := w.prev[f]
		results := diff.Compare(map[string]map[string]string{"prev": prev, "curr": curr})
		if len(results) > 0 {
			w.Events <- Event{File: f, At: time.Now(), Results: results}
		}
		w.prev[f] = curr
	}
}
