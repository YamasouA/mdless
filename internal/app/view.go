package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/YamasouA/mdless/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if len(m.Tabs) == 0 {
		return ""
	}

	header := m.renderTabs()
	body := m.renderContent()
	footer := m.renderStatus()

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m Model) renderTabs() string {
	parts := make([]string, 0, len(m.Tabs)*2)
	for i, tab := range m.Tabs {
		if i > 0 {
			parts = append(parts, ui.TabSeparator.Render(" "))
		}
		label := fmt.Sprintf("%d %s", i+1, truncateTabName(tabName(tab.Page.Path), tabNameLimit(m.Width, len(m.Tabs))))
		if i == m.CurrentTab {
			parts = append(parts, ui.ActiveTab.Render(label))
			continue
		}
		parts = append(parts, ui.InactiveTab.Render(label))
	}
	return ui.Header.Width(m.Width).Render(lipgloss.JoinHorizontal(lipgloss.Top, parts...))
}

func tabName(path string) string {
	name := filepath.Base(path)
	if name == "." {
		return path
	}
	return name
}

func tabNameLimit(width, tabs int) int {
	if tabs <= 0 {
		return 16
	}
	if width <= 0 {
		return 18
	}
	limit := width/tabs - 6
	if limit < 8 {
		return 8
	}
	if limit > 24 {
		return 24
	}
	return limit
}

func truncateTabName(name string, limit int) string {
	runes := []rune(name)
	if len(runes) <= limit {
		return name
	}
	if limit <= 1 {
		return string(runes[:limit])
	}
	return string(runes[:limit-1]) + "…"
}

func (m Model) renderContent() string {
	tab := m.currentTab()
	height := m.viewportHeight()
	start := tab.ScrollY
	end := start + height
	if end > len(tab.Page.Content) {
		end = len(tab.Page.Content)
	}

	lines := append([]string(nil), tab.Page.Content[start:end]...)
	if m.SearchQuery != "" && len(m.Matches) > 0 {
		for i, line := range lines {
			lines[i] = highlightMatches(line, m.SearchQuery)
		}
	}
	for len(lines) < height {
		lines = append(lines, "")
	}

	if m.Mode == ModeLinks {
		return m.renderLinks(height)
	}

	return strings.Join(lines, "\n")
}

func (m Model) renderLinks(height int) string {
	links := m.currentTab().Page.Links
	lines := make([]string, 0, len(links))
	for i, link := range links {
		prefix := "  "
		if i == m.LinkIndex {
			prefix = "> "
		}
		lines = append(lines, fmt.Sprintf("%s%d. %s -> %s", prefix, i+1, link.Text, link.Target))
	}
	if len(lines) > height {
		lines = lines[:height]
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderStatus() string {
	tab := m.currentTab()
	status := fmt.Sprintf("%s | %d/%d | tab %d/%d", tab.Page.Path, tab.ScrollY+1, m.maxScroll()+1, m.CurrentTab+1, len(m.Tabs))
	if m.Mode == ModeSearch {
		status = "/" + m.SearchInput
	}
	if m.Status != "" && m.Mode != ModeSearch {
		status += " | " + m.Status
	}
	return ui.Footer.Width(m.Width).Render(status)
}
