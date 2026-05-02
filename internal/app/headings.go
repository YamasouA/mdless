package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) openHeadingList() {
	if len(m.currentTab().Page.Headings) == 0 {
		m.Status = "no headings"
		return
	}
	m.Mode = ModeHeadings
	m.HeadingIndex = m.nearestHeadingIndex()
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
