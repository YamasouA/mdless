package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

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
