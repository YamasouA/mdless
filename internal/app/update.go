package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/YamasouA/mdview/internal/nav"
	"github.com/YamasouA/mdview/internal/render"
	"github.com/YamasouA/mdview/internal/watch"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Ready = true
		m.clampScroll()
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	case watch.Event:
		return m.handleWatchEvent(msg)
	case watch.Error:
		return m.handleWatchError(msg), nil
	default:
		return m, nil
	}
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyCtrlC || msg.String() == "q" && m.Mode == ModeView {
		return m, tea.Quit
	}

	switch m.Mode {
	case ModeSearch:
		return m.handleSearchKey(msg), nil
	case ModeLinks:
		return m.handleLinksKey(msg)
	case ModeHeadings:
		return m.handleHeadingsKey(msg), nil
	default:
		return m.handleViewKey(msg)
	}
}

func (m Model) handleViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.PendingHeadingPrefix != "" {
		prefix := m.PendingHeadingPrefix
		m.PendingHeadingPrefix = ""
		if msg.String() == "h" {
			if prefix == "]" {
				m.moveHeading(1)
			} else {
				m.moveHeading(-1)
			}
		}
		return m, nil
	}

	if m.PendingG {
		m.PendingG = false
		switch msg.String() {
		case "t":
			m.switchTab(1)
			return m, nil
		case "T":
			m.switchTab(-1)
			return m, nil
		case "g":
			m.scrollTop()
			return m, nil
		default:
			return m, nil
		}
	}

	switch msg.String() {
	case "j":
		m.scrollBy(1)
	case "k":
		m.scrollBy(-1)
	case "ctrl+d":
		m.scrollBy(halfPage(m.viewportHeight()))
	case "ctrl+u":
		m.scrollBy(-halfPage(m.viewportHeight()))
	case "g":
		m.PendingG = true
	case "G":
		tab := m.currentTab()
		tab.ScrollY = m.maxScroll()
		m.setCurrentTab(tab)
	case "/":
		m.Mode = ModeSearch
		m.SearchInput = ""
		m.Status = "search"
	case "n":
		m.moveMatch(1)
	case "N":
		m.moveMatch(-1)
	case "H":
		m.openHeadingList()
	case "]", "[":
		m.PendingHeadingPrefix = msg.String()
	case "o":
		m.openLinkList()
	case "enter":
		return m, m.openLink(0, false)
	case "t":
		return m, m.openLink(0, true)
	case "b":
		return m, m.goHistory(-1)
	case "f":
		return m, m.goHistory(1)
	case "x":
		m.closeCurrentTab()
	}
	return m, nil
}

func (m *Model) openHeadingList() {
	if len(m.currentTab().Page.Headings) == 0 {
		m.Status = "no headings"
		return
	}
	m.Mode = ModeHeadings
	m.HeadingIndex = m.nearestHeadingIndex()
}

func (m *Model) openLinkList() {
	if len(m.currentTab().Page.Links) == 0 {
		m.Status = "no links"
		return
	}
	m.Mode = ModeLinks
	m.LinkIndex = 0
}

func (m Model) handleSearchKey(msg tea.KeyMsg) Model {
	switch msg.Type {
	case tea.KeyEsc:
		m.Mode = ModeView
		m.SearchInput = ""
		m.Status = ""
	case tea.KeyEnter:
		m.Mode = ModeView
		m.SearchQuery = m.SearchInput
		m.performSearch()
	case tea.KeyBackspace:
		if len(m.SearchInput) > 0 {
			m.SearchInput = m.SearchInput[:len(m.SearchInput)-1]
		}
	case tea.KeyRunes:
		m.SearchInput += string(msg.Runes)
	}
	return m
}

func (m Model) handleLinksKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	links := m.currentTab().Page.Links
	switch msg.Type {
	case tea.KeyEsc:
		m.Mode = ModeView
	case tea.KeyEnter:
		return m, m.openLink(m.LinkIndex, false)
	case tea.KeyRunes:
		key := string(msg.Runes)
		switch key {
		case "j":
			if m.LinkIndex < len(links)-1 {
				m.LinkIndex++
			}
		case "k":
			if m.LinkIndex > 0 {
				m.LinkIndex--
			}
		default:
			if len(key) == 1 && key[0] >= '1' && key[0] <= '9' {
				return m, m.openLink(int(key[0]-'1'), false)
			}
			if key == "t" {
				return m, m.openLink(m.LinkIndex, true)
			}
		}
	}
	return m, nil
}

func (m Model) handleHeadingsKey(msg tea.KeyMsg) Model {
	headings := m.currentTab().Page.Headings
	switch msg.Type {
	case tea.KeyEsc:
		m.Mode = ModeView
	case tea.KeyEnter:
		m.jumpToHeading(m.HeadingIndex)
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "j":
			if m.HeadingIndex < len(headings)-1 {
				m.HeadingIndex++
			}
		case "k":
			if m.HeadingIndex > 0 {
				m.HeadingIndex--
			}
		}
	}
	return m
}

func (m *Model) scrollBy(delta int) {
	tab := m.currentTab()
	tab.ScrollY += delta
	m.setCurrentTab(tab)
	m.clampScroll()
}

func (m *Model) scrollTop() {
	tab := m.currentTab()
	tab.ScrollY = 0
	m.setCurrentTab(tab)
}

func (m *Model) clampScroll() {
	tab := m.currentTab()
	if tab.ScrollY < 0 {
		tab.ScrollY = 0
	}
	if tab.ScrollY > m.maxScroll() {
		tab.ScrollY = m.maxScroll()
	}
	m.setCurrentTab(tab)
}

func (m *Model) performSearch() {
	m.refreshSearchMatches()
	if m.SearchQuery == "" {
		m.Status = ""
		return
	}
	if len(m.Matches) == 0 {
		m.Status = "no matches"
		return
	}
	m.scrollToLine(m.Matches[0].Line)
	m.Status = fmt.Sprintf("%d/%d", 1, len(m.Matches))
}

func (m *Model) refreshSearchMatches() {
	m.Matches = nil
	m.MatchIndex = 0
	if m.SearchQuery == "" {
		return
	}
	for i, line := range m.currentTab().Page.Content {
		m.Matches = append(m.Matches, findPlainMatches(i, line, m.SearchQuery)...)
	}
}

func (m *Model) moveMatch(delta int) {
	if len(m.Matches) == 0 {
		if m.SearchQuery != "" {
			m.performSearch()
		}
		return
	}
	m.MatchIndex = (m.MatchIndex + delta + len(m.Matches)) % len(m.Matches)
	m.scrollToLine(m.Matches[m.MatchIndex].Line)
	m.Status = fmt.Sprintf("%d/%d", m.MatchIndex+1, len(m.Matches))
}

func (m *Model) moveHeading(delta int) {
	headings := m.currentTab().Page.Headings
	if len(headings) == 0 {
		m.Status = "no headings"
		return
	}

	current := m.currentTab().ScrollY
	target := 0
	if delta > 0 {
		target = len(headings) - 1
		for i, heading := range headings {
			if heading.Line > current {
				target = i
				break
			}
		}
	} else {
		for i := len(headings) - 1; i >= 0; i-- {
			if headings[i].Line < current {
				target = i
				break
			}
		}
	}
	m.jumpToHeading(target)
}

func (m *Model) jumpToHeading(index int) {
	headings := m.currentTab().Page.Headings
	if index < 0 || index >= len(headings) {
		m.Status = "no headings"
		return
	}
	m.HeadingIndex = index
	m.scrollToLine(headings[index].Line)
	m.Mode = ModeView
	m.Status = fmt.Sprintf("heading: %s", headings[index].Text)
}

func (m Model) nearestHeadingIndex() int {
	headings := m.currentTab().Page.Headings
	if len(headings) == 0 {
		return 0
	}
	current := m.currentTab().ScrollY
	index := 0
	for i, heading := range headings {
		if heading.Line > current {
			break
		}
		index = i
	}
	return index
}

func (m *Model) scrollToLine(line int) {
	tab := m.currentTab()
	tab.ScrollY = line
	m.setCurrentTab(tab)
	m.clampScroll()
}

func (m *Model) openLink(index int, newTab bool) tea.Cmd {
	tab := m.currentTab()
	if index < 0 || index >= len(tab.Page.Links) {
		if len(tab.Page.Links) == 0 {
			m.Status = "no links"
		} else {
			m.Status = "no link"
		}
		m.Mode = ModeView
		return nil
	}

	link := tab.Page.Links[index]
	if strings.Contains(link.Target, "://") {
		m.Status = "external links are not supported"
		m.Mode = ModeView
		return nil
	}

	target := render.ResolveTarget(tab.Page.Path, link.Target)
	if newTab {
		if index := m.findTabByPath(target); index >= 0 {
			m.switchToTab(index, "switched to existing tab")
			return nil
		}
	}

	page, err := m.LoadPage(target)
	if err != nil {
		m.Status = err.Error()
		m.Mode = ModeView
		return nil
	}

	if newTab {
		m.Tabs = append(m.Tabs, Tab{
			Page:    page,
			History: nav.NewHistory(nav.HistoryEntry{Path: page.Path}),
		})
		m.CurrentTab = len(m.Tabs) - 1
	} else {
		tab.Page = page
		tab.ScrollY = 0
		tab.History = tab.History.UpdateCurrent(m.currentTab().ScrollY)
		tab.History = tab.History.Push(nav.HistoryEntry{Path: page.Path})
		m.setCurrentTab(tab)
	}
	m.Mode = ModeView
	m.Status = ""
	m.PendingG = false
	m.SearchInput = ""
	m.SearchQuery = ""
	m.Matches = nil
	m.MatchIndex = 0
	return m.watchPathCmd(page.Path)
}

func (m *Model) findTabByPath(path string) int {
	cleanPath := filepath.Clean(path)
	for i, tab := range m.Tabs {
		if filepath.Clean(tab.Page.Path) == cleanPath {
			return i
		}
	}
	return -1
}

func (m *Model) switchToTab(index int, status string) {
	if index < 0 || index >= len(m.Tabs) {
		return
	}
	query := m.SearchQuery
	m.CurrentTab = index
	m.resetTransientState()
	m.SearchQuery = query
	m.refreshSearchMatches()
	m.Status = status
}

func (m *Model) goHistory(direction int) tea.Cmd {
	tab := m.currentTab()
	tab.History = tab.History.UpdateCurrent(tab.ScrollY)
	oldHistory := tab.History

	var entry nav.HistoryEntry
	var ok bool
	if direction < 0 {
		tab.History, entry, ok = tab.History.Back()
	} else {
		tab.History, entry, ok = tab.History.Forward()
	}
	if !ok {
		m.Status = "no history"
		return nil
	}

	page, err := m.LoadPage(entry.Path)
	if err != nil {
		tab.History = oldHistory
		m.setCurrentTab(tab)
		m.Status = err.Error()
		return nil
	}

	tab.Page = page
	tab.ScrollY = entry.ScrollY
	m.setCurrentTab(tab)
	m.clampScroll()
	m.Status = ""
	return m.watchPathCmd(page.Path)
}

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

func (m *Model) switchTab(delta int) {
	if len(m.Tabs) <= 1 {
		m.Status = "no other tabs"
		return
	}
	m.switchToTab((m.CurrentTab+delta+len(m.Tabs))%len(m.Tabs), "")
}

func (m *Model) closeCurrentTab() {
	if len(m.Tabs) <= 1 {
		m.Status = "cannot close last tab"
		return
	}

	m.Tabs = append(m.Tabs[:m.CurrentTab], m.Tabs[m.CurrentTab+1:]...)
	if m.CurrentTab >= len(m.Tabs) {
		m.CurrentTab = len(m.Tabs) - 1
	}
	query := m.SearchQuery
	m.resetTransientState()
	m.SearchQuery = query
	m.refreshSearchMatches()
}

func (m *Model) resetTransientState() {
	m.Mode = ModeView
	m.Status = ""
	m.SearchInput = ""
	m.SearchQuery = ""
	m.Matches = nil
	m.MatchIndex = 0
	m.LinkIndex = 0
	m.HeadingIndex = 0
	m.PendingG = false
	m.PendingHeadingPrefix = ""
}

func halfPage(height int) int {
	if height <= 1 {
		return 1
	}
	return height / 2
}

func clampScrollY(scrollY, contentLines, viewportHeight int) int {
	if scrollY < 0 {
		return 0
	}
	maxScroll := contentLines - viewportHeight
	if maxScroll < 0 {
		maxScroll = 0
	}
	if scrollY > maxScroll {
		return maxScroll
	}
	return scrollY
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
