package app

import (
	"github.com/YamasouA/mdview/internal/nav"
	"github.com/YamasouA/mdview/internal/render"
	"github.com/YamasouA/mdview/internal/watch"
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	ModeView Mode = iota
	ModeSearch
	ModeLinks
)

type PageLoader func(path string) (render.Page, error)

type Model struct {
	Tabs       []Tab
	CurrentTab int
	Width      int
	Height     int
	Ready      bool

	Mode        Mode
	Status      string
	SearchInput string
	SearchQuery string
	Matches     []int
	MatchIndex  int
	LinkIndex   int
	PendingG    bool

	LoadPage     PageLoader
	LiveReload   bool
	WatchedPaths map[string]bool
}

type Tab struct {
	Page    render.Page
	History nav.History
	ScrollY int
}

func NewModel(page render.Page, loader PageLoader) Model {
	return NewModelWithPages([]render.Page{page}, loader)
}

func NewModelWithPages(pages []render.Page, loader PageLoader) Model {
	tabs := make([]Tab, 0, len(pages))
	for _, page := range pages {
		tabs = append(tabs, Tab{
			Page:    page,
			History: nav.NewHistory(nav.HistoryEntry{Path: page.Path}),
		})
	}
	return Model{
		Tabs:         tabs,
		LoadPage:     loader,
		LiveReload:   true,
		WatchedPaths: make(map[string]bool),
	}
}

func (m Model) Init() tea.Cmd {
	return m.watchOpenPathsCmd()
}

func (m Model) currentTab() Tab {
	if len(m.Tabs) == 0 || m.CurrentTab < 0 || m.CurrentTab >= len(m.Tabs) {
		return Tab{}
	}
	return m.Tabs[m.CurrentTab]
}

func (m *Model) setCurrentTab(tab Tab) {
	if len(m.Tabs) == 0 || m.CurrentTab < 0 || m.CurrentTab >= len(m.Tabs) {
		return
	}
	m.Tabs[m.CurrentTab] = tab
}

func (m Model) viewportHeight() int {
	if m.Height <= 2 {
		return 1
	}
	return m.Height - 2
}

func (m Model) maxScroll() int {
	maxScroll := len(m.currentTab().Page.Content) - m.viewportHeight()
	if maxScroll < 0 {
		return 0
	}
	return maxScroll
}

func (m *Model) watchPathCmd(path string) tea.Cmd {
	if !m.LiveReload || path == "" {
		return nil
	}
	if m.WatchedPaths == nil {
		m.WatchedPaths = make(map[string]bool)
	}
	if m.WatchedPaths[path] {
		return nil
	}
	m.WatchedPaths[path] = true
	return watch.Watch(path)
}

func (m *Model) watchOpenPathsCmd() tea.Cmd {
	if !m.LiveReload {
		return nil
	}
	var cmds []tea.Cmd
	for _, tab := range m.Tabs {
		if cmd := m.watchPathCmd(tab.Page.Path); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}
