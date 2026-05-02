package app

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

func (m *Model) scrollToLine(line int) {
	tab := m.currentTab()
	tab.ScrollY = line
	m.setCurrentTab(tab)
	m.clampScroll()
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
