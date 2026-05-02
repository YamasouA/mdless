package app

import (
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
