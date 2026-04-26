package watch

import (
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fsnotify/fsnotify"
)

const DefaultDebounce = 150 * time.Millisecond

type Event struct {
	Path string
}

type Error struct {
	Path string
	Err  error
}

func Watch(path string) tea.Cmd {
	return func() tea.Msg {
		cleanPath := filepath.Clean(path)
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return Error{Path: path, Err: err}
		}
		defer watcher.Close()

		if err := watcher.Add(filepath.Dir(cleanPath)); err != nil {
			return Error{Path: path, Err: err}
		}

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return Error{Path: path, Err: fsnotify.ErrEventOverflow}
				}
				if filepath.Clean(event.Name) != cleanPath || !isReloadEvent(event) {
					continue
				}
				time.Sleep(DefaultDebounce)
				return Event{Path: path}
			case err, ok := <-watcher.Errors:
				if !ok {
					return Error{Path: path, Err: fsnotify.ErrEventOverflow}
				}
				return Error{Path: path, Err: err}
			}
		}
	}
}

func isReloadEvent(event fsnotify.Event) bool {
	return event.Has(fsnotify.Write) ||
		event.Has(fsnotify.Create) ||
		event.Has(fsnotify.Rename) ||
		event.Has(fsnotify.Remove)
}
