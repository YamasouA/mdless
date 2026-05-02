package app

import (
	"strings"

	"github.com/YamasouA/mdview/internal/watch"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleWatchEvent(event watch.Event) (tea.Model, tea.Cmd) {
	if m.WatchedPaths != nil {
		m.WatchedPaths[event.Path] = false
	}

	if !m.hasOpenPath(event.Path) {
		return m, nil
	}

	page, err := m.LoadPage(event.Path)
	if err != nil {
		m.Status = err.Error()
		return m, m.watchPathCmd(event.Path)
	}

	for i, tab := range m.Tabs {
		if tab.Page.Path != event.Path {
			continue
		}
		anchor := scrollAnchor(tab, m.viewportHeight())
		tab.Page = page
		tab.ScrollY = restoreScrollY(tab.ScrollY, anchor, tab.Page.Content, m.viewportHeight())
		m.Tabs[i] = tab
	}
	if m.SearchQuery != "" {
		m.refreshSearchMatches()
	}

	m.Status = "reloaded"
	return m, m.watchPathCmd(event.Path)
}

func (m Model) handleWatchError(event watch.Error) Model {
	if m.WatchedPaths != nil {
		m.WatchedPaths[event.Path] = false
	}
	if event.Err != nil {
		m.Status = event.Err.Error()
	}
	return m
}

func (m Model) hasOpenPath(path string) bool {
	for _, tab := range m.Tabs {
		if tab.Page.Path == path {
			return true
		}
	}
	return false
}

func scrollAnchor(tab Tab, viewportHeight int) string {
	if len(tab.Page.Content) == 0 {
		return ""
	}
	start := tab.ScrollY
	if start < 0 {
		start = 0
	}
	if start >= len(tab.Page.Content) {
		start = len(tab.Page.Content) - 1
	}
	end := start + viewportHeight
	if end > len(tab.Page.Content) {
		end = len(tab.Page.Content)
	}
	for _, line := range tab.Page.Content[start:end] {
		text := strings.TrimSpace(plainText(line))
		if text != "" {
			return text
		}
	}
	return strings.TrimSpace(plainText(tab.Page.Content[start]))
}

func restoreScrollY(previous int, anchor string, content []string, viewportHeight int) int {
	if anchor != "" {
		found := -1
		for i, line := range content {
			if strings.TrimSpace(plainText(line)) == anchor {
				if found >= 0 {
					return clampScrollY(previous, len(content), viewportHeight)
				}
				found = i
			}
		}
		if found >= 0 {
			return clampScrollY(found, len(content), viewportHeight)
		}
	}
	return clampScrollY(previous, len(content), viewportHeight)
}
