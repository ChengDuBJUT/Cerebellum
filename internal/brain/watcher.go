package brain

import (
	"log"
	"os"
	"time"

	"cerebellum/internal/store"
	"github.com/fsnotify/fsnotify"
)

// Watcher monitors brain.md for changes
type Watcher struct {
	path     string
	watcher  *fsnotify.Watcher
	store    *store.MarkdownStore
	pollStop chan struct{}
	pollDone chan struct{}
}

// NewWatcher creates a new file watcher
func NewWatcher(path string, store *store.MarkdownStore) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// Ensure file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	w := &Watcher{
		path:     path,
		watcher:  watcher,
		store:    store,
		pollStop: make(chan struct{}),
		pollDone: make(chan struct{}),
	}

	return w, nil
}

// Start begins watching the file
func (w *Watcher) Start() {
	// Add file to watcher
	if err := w.watcher.Add(w.path); err != nil {
		log.Printf("Warning: Could not add watcher for %s: %v", w.path, err)
	}

	// Also watch the directory
	dir := w.path
	for i := len(dir) - 1; i >= 0; i-- {
		if dir[i] == '/' || dir[i] == '\\' {
			dir = dir[:i]
			break
		}
	}
	if dir != "" {
		w.watcher.Add(dir)
	}

	// Start file watcher goroutine
	go w.watchEvents()

	// Start polling goroutine (fallback for Windows)
	go w.pollFile()

	log.Printf("Started watching %s for changes", w.path)
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	close(w.pollStop)
	<-w.pollDone
	w.watcher.Close()
	log.Println("File watcher stopped")
}

func (w *Watcher) watchEvents() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Printf("Detected change in %s", event.Name)
				w.handleChange()
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (w *Watcher) pollFile() {
	defer close(w.pollDone)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if w.store.HasChanged() {
				log.Printf("Polling detected change in %s", w.path)
				w.handleChange()
			}
		case <-w.pollStop:
			return
		}
	}
}

func (w *Watcher) handleChange() {
	// Reload the markdown store
	if err := w.store.Reload(); err != nil {
		log.Printf("Error reloading brain.md: %v", err)
		return
	}

	log.Printf("Reloaded %d tasks from brain.md", len(w.store.GetTasks()))
}
