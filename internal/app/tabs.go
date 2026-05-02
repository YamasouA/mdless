package app

func (m *Model) switchTab(delta int) {
	if len(m.Tabs) <= 1 {
		m.Status = "no other tabs"
		return
	}
	m.switchToTab((m.CurrentTab+delta+len(m.Tabs))%len(m.Tabs), "")
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
