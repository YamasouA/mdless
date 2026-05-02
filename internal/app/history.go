package app

import (
	"github.com/YamasouA/mdview/internal/nav"
	tea "github.com/charmbracelet/bubbletea"
)

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
