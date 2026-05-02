package app

import (
	"path/filepath"
	"strings"

	"github.com/YamasouA/mdview/internal/nav"
	"github.com/YamasouA/mdview/internal/render"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) openLinkList() {
	if len(m.currentTab().Page.Links) == 0 {
		m.Status = "no links"
		return
	}
	m.Mode = ModeLinks
	m.LinkIndex = 0
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
