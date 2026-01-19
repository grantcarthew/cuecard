package watcher

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Event represents a file system event
type Event struct {
	Path string
	Op   Operation
}

// Operation represents the type of file system operation
type Operation int

const (
	Create Operation = iota
	Write
	Remove
	Rename
)

// Watcher watches a directory for file changes
type Watcher struct {
	fsWatcher *fsnotify.Watcher
	dir       string
	onChange  func()
	done      chan struct{}
	mu        sync.Mutex
	debounce  time.Duration
}

// New creates a new file watcher for the given directory
func New(dir string, onChange func()) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		fsWatcher: fsWatcher,
		dir:       dir,
		onChange:  onChange,
		done:      make(chan struct{}),
		debounce:  100 * time.Millisecond,
	}

	return w, nil
}

// Start begins watching the directory
func (w *Watcher) Start() error {
	if err := w.fsWatcher.Add(w.dir); err != nil {
		return err
	}

	go w.watch()
	return nil
}

// Stop stops watching the directory
func (w *Watcher) Stop() error {
	close(w.done)
	return w.fsWatcher.Close()
}

func (w *Watcher) watch() {
	var timer *time.Timer
	var timerMu sync.Mutex

	for {
		select {
		case <-w.done:
			return
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}

			// Only watch markdown files
			if !isMarkdownFile(event.Name) {
				continue
			}

			// Debounce events
			timerMu.Lock()
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(w.debounce, func() {
				w.mu.Lock()
				defer w.mu.Unlock()
				if w.onChange != nil {
					w.onChange()
				}
			})
			timerMu.Unlock()

		case _, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			// Log errors but continue watching
		}
	}
}

func isMarkdownFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md"
}

// SetDebounce sets the debounce duration for events
func (w *Watcher) SetDebounce(d time.Duration) {
	w.debounce = d
}
