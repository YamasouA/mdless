package app

import (
	"errors"
	"strings"
	"testing"

	"github.com/YamasouA/mdview/internal/nav"
	"github.com/YamasouA/mdview/internal/render"
	"github.com/YamasouA/mdview/internal/watch"
	tea "github.com/charmbracelet/bubbletea"
)

func TestScrollClamps(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: lines(10)})
	m.Height = 5

	m = press(m, "G")
	if got, want := m.currentTab().ScrollY, 7; got != want {
		t.Fatalf("ScrollY = %d, want %d", got, want)
	}

	m = press(m, "j")
	if got, want := m.currentTab().ScrollY, 7; got != want {
		t.Fatalf("ScrollY after bottom j = %d, want %d", got, want)
	}

	m = pressSequence(m, "g", "g")
	m = press(m, "k")
	if got := m.currentTab().ScrollY; got != 0 {
		t.Fatalf("ScrollY after top k = %d, want 0", got)
	}
}

func TestGGScrollsTop(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: lines(10)})
	m.Height = 5
	m.Tabs[0].ScrollY = 3

	m = press(m, "g")
	if got := m.currentTab().ScrollY; got != 3 {
		t.Fatalf("ScrollY after pending g = %d, want 3", got)
	}

	m = press(m, "g")
	if got := m.currentTab().ScrollY; got != 0 {
		t.Fatalf("ScrollY after gg = %d, want 0", got)
	}
}

func TestUnknownGSequenceDoesNothing(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: lines(10)})
	m.Height = 5
	m.Tabs[0].ScrollY = 3

	m = pressSequence(m, "g", "j")
	if got := m.currentTab().ScrollY; got != 3 {
		t.Fatalf("ScrollY after unknown g sequence = %d, want 3", got)
	}
}

func TestSearchMovesToMatch(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: []string{"alpha", "beta", "target", "gamma", "target"}})
	m.Height = 4

	m = press(m, "/")
	m = typeRunes(m, "target")
	m = key(m, tea.KeyEnter)

	if got := m.currentTab().ScrollY; got != 2 {
		t.Fatalf("ScrollY = %d, want 2", got)
	}
	if got := m.Status; got != "1/2" {
		t.Fatalf("Status = %q, want 1/2", got)
	}

	m = press(m, "n")
	if got := m.Status; got != "2/2" {
		t.Fatalf("Status after n = %q, want 2/2", got)
	}
}

func TestSearchMatchesANSIStyledContent(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: []string{"alpha", "\x1b[31mtarget\x1b[0m"}})
	m.Height = 3

	m = press(m, "/")
	m = typeRunes(m, "target")
	m = key(m, tea.KeyEnter)

	if got := m.currentTab().ScrollY; got != 1 {
		t.Fatalf("ScrollY = %d, want 1", got)
	}
	if got := m.Status; got != "1/1" {
		t.Fatalf("Status = %q, want 1/1", got)
	}
}

func TestRenderContentHighlightsSearchMatches(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: []string{"target and Target"}})
	m.Height = 3
	m.SearchQuery = "target"
	m.Matches = []int{0}

	body := m.renderContent()

	if count := strings.Count(body, highlightStart); count != 2 {
		t.Fatalf("highlight count = %d, want 2 in %q", count, body)
	}
}

func TestSearchNoMatchKeepsScroll(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: lines(10)})
	m.Height = 5
	m = press(m, "j")

	m = press(m, "/")
	m = typeRunes(m, "missing")
	m = key(m, tea.KeyEnter)

	if got := m.currentTab().ScrollY; got != 1 {
		t.Fatalf("ScrollY = %d, want 1", got)
	}
	if got := m.Status; got != "no matches" {
		t.Fatalf("Status = %q, want no matches", got)
	}
}

func TestOpenLinkUpdatesPageAndHistory(t *testing.T) {
	next := render.Page{Path: "next.md", Content: []string{"next"}}
	m := NewModel(render.Page{
		Path:    "a.md",
		Content: []string{"home"},
		Links:   []nav.Link{{Text: "Next", Target: "next.md"}},
	}, func(path string) (render.Page, error) {
		if path != "next.md" {
			t.Fatalf("path = %q, want next.md", path)
		}
		return next, nil
	})

	m = press(m, "enter")

	if got := m.currentTab().Page.Path; got != "next.md" {
		t.Fatalf("Path = %q, want next.md", got)
	}
	if got := len(m.currentTab().History.Entries); got != 2 {
		t.Fatalf("history len = %d, want 2", got)
	}
}

func TestOpenLinkKeepsOtherInitialTabs(t *testing.T) {
	pages := map[string]render.Page{
		"README.md": {
			Path:    "README.md",
			Content: []string{"readme"},
			Links:   []nav.Link{{Text: "Next", Target: "docs/next.md"}},
		},
		"README.en.md": {Path: "README.en.md", Content: []string{"english"}},
		"docs/next.md": {Path: "docs/next.md", Content: []string{"next"}},
	}
	m := NewModelWithPages([]render.Page{pages["README.md"], pages["README.en.md"]}, func(path string) (render.Page, error) {
		return pages[path], nil
	})
	m.Width = 80

	m = press(m, "enter")

	if got := len(m.Tabs); got != 2 {
		t.Fatalf("len(Tabs) = %d, want 2", got)
	}
	if got := m.Tabs[0].Page.Path; got != "docs/next.md" {
		t.Fatalf("Tabs[0].Path = %q, want docs/next.md", got)
	}
	if got := m.Tabs[1].Page.Path; got != "README.en.md" {
		t.Fatalf("Tabs[1].Path = %q, want README.en.md", got)
	}

	view := m.View()
	for _, text := range []string{"1 next.md", "2 README.en.md", "tab 1/2"} {
		if !strings.Contains(view, text) {
			t.Fatalf("View() missing %q after opening link: %q", text, view)
		}
	}
}

func TestOpenLinkSavesCurrentScrollInHistory(t *testing.T) {
	pages := map[string]render.Page{
		"a.md": {
			Path:    "a.md",
			Content: lines(10),
			Links:   []nav.Link{{Text: "Next", Target: "next.md"}},
		},
		"next.md": {Path: "next.md", Content: []string{"next"}},
	}
	m := NewModel(pages["a.md"], func(path string) (render.Page, error) {
		return pages[path], nil
	})
	m.Height = 5
	m = press(m, "j")
	m = press(m, "j")

	m = press(m, "enter")
	m = press(m, "b")

	if got := m.currentTab().ScrollY; got != 2 {
		t.Fatalf("ScrollY after back = %d, want 2", got)
	}
}

func TestHistoryNavigationSavesCurrentScroll(t *testing.T) {
	pages := map[string]render.Page{
		"a.md": {Path: "a.md", Content: []string{"a"}},
		"b.md": {Path: "b.md", Content: lines(10)},
	}
	m := NewModel(pages["a.md"], func(path string) (render.Page, error) {
		return pages[path], nil
	})
	m.Height = 5
	tab := m.currentTab()
	tab.History = tab.History.Push(nav.HistoryEntry{Path: "b.md", ScrollY: 0})
	tab.Page = pages["b.md"]
	tab.ScrollY = 3
	m.setCurrentTab(tab)

	m = press(m, "b")
	m = press(m, "f")

	if got := m.currentTab().ScrollY; got != 3 {
		t.Fatalf("ScrollY after forward = %d, want 3", got)
	}
}

func TestOpenLinkFailureKeepsPage(t *testing.T) {
	m := NewModel(render.Page{
		Path:    "a.md",
		Content: []string{"home"},
		Links:   []nav.Link{{Text: "Missing", Target: "missing.md"}},
	}, func(string) (render.Page, error) {
		return render.Page{}, errors.New("missing")
	})

	m = press(m, "enter")

	if got := m.currentTab().Page.Path; got != "a.md" {
		t.Fatalf("Path = %q, want a.md", got)
	}
	if got := m.Status; got != "missing" {
		t.Fatalf("Status = %q, want missing", got)
	}
}

func TestHistoryBackForwardLoadsSavedPages(t *testing.T) {
	pages := map[string]render.Page{
		"a.md": {Path: "a.md", Content: []string{"a"}},
		"b.md": {Path: "b.md", Content: []string{"b"}},
	}
	m := NewModel(pages["a.md"], func(path string) (render.Page, error) {
		return pages[path], nil
	})
	tab := m.currentTab()
	tab.History = tab.History.Push(nav.HistoryEntry{Path: "b.md", ScrollY: 0})
	tab.Page = pages["b.md"]
	m.setCurrentTab(tab)

	m = press(m, "b")
	if got := m.currentTab().Page.Path; got != "a.md" {
		t.Fatalf("Path after back = %q, want a.md", got)
	}

	m = press(m, "f")
	if got := m.currentTab().Page.Path; got != "b.md" {
		t.Fatalf("Path after forward = %q, want b.md", got)
	}
}

func TestViewShowsEnterLinkTarget(t *testing.T) {
	m := NewModel(render.Page{
		Path:    "a.md",
		Content: []string{"home"},
		Links:   []nav.Link{{Text: "Next", Target: "next.md"}},
	}, func(path string) (render.Page, error) {
		return render.Page{Path: path, Content: []string{path}}, nil
	})

	view := m.View()

	if !strings.Contains(view, "enter: Next -> next.md") {
		t.Fatalf("View() does not show enter link target: %q", view)
	}
}

func TestViewKeepsEnterLinkTargetAfterReloadStatus(t *testing.T) {
	m := NewModel(render.Page{
		Path:    "a.md",
		Content: []string{"home"},
		Links:   []nav.Link{{Text: "Next", Target: "next.md"}},
	}, func(path string) (render.Page, error) {
		return render.Page{Path: path, Content: []string{path}}, nil
	})
	m.Status = "reloaded"

	view := m.View()

	if !strings.Contains(view, "reloaded") {
		t.Fatalf("View() does not show reload status: %q", view)
	}
	if !strings.Contains(view, "enter: Next -> next.md") {
		t.Fatalf("View() does not keep enter link target: %q", view)
	}
}

func TestEnterWithNoLinksShowsStatus(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: []string{"home"}})

	m = press(m, "enter")

	if got := m.Mode; got != ModeView {
		t.Fatalf("Mode = %v, want ModeView", got)
	}
	if got := m.Status; got != "no links" {
		t.Fatalf("Status = %q, want no links", got)
	}
}

func TestOOpensLinkListInViewMode(t *testing.T) {
	m := NewModel(render.Page{
		Path:    "a.md",
		Content: []string{"home"},
		Links:   []nav.Link{{Text: "Next", Target: "next.md"}},
	}, func(path string) (render.Page, error) {
		return render.Page{Path: path, Content: []string{path}}, nil
	})

	m = press(m, "o")

	if got := m.Mode; got != ModeLinks {
		t.Fatalf("Mode = %v, want ModeLinks", got)
	}
	if got := m.currentTab().Page.Path; got != "a.md" {
		t.Fatalf("Path = %q, want a.md", got)
	}
}

func TestOpenLinkInNewTab(t *testing.T) {
	pages := map[string]render.Page{
		"a.md": {
			Path:    "a.md",
			Content: []string{"home"},
			Links:   []nav.Link{{Text: "Next", Target: "next.md"}},
		},
		"next.md": {Path: "next.md", Content: []string{"next"}},
	}
	m := NewModel(pages["a.md"], func(path string) (render.Page, error) {
		return pages[path], nil
	})

	m = press(m, "t")

	if got := len(m.Tabs); got != 2 {
		t.Fatalf("len(Tabs) = %d, want 2", got)
	}
	if got := m.CurrentTab; got != 1 {
		t.Fatalf("CurrentTab = %d, want 1", got)
	}
	if got := m.currentTab().Page.Path; got != "next.md" {
		t.Fatalf("current path = %q, want next.md", got)
	}
	if got := m.Tabs[0].Page.Path; got != "a.md" {
		t.Fatalf("first tab path = %q, want a.md", got)
	}
}

func TestOpenLinkInNewTabFromLinkList(t *testing.T) {
	pages := map[string]render.Page{
		"a.md": {
			Path:    "a.md",
			Content: []string{"home"},
			Links: []nav.Link{
				{Text: "One", Target: "one.md"},
				{Text: "Two", Target: "two.md"},
			},
		},
		"one.md": {Path: "one.md", Content: []string{"one"}},
		"two.md": {Path: "two.md", Content: []string{"two"}},
	}
	m := NewModel(pages["a.md"], func(path string) (render.Page, error) {
		return pages[path], nil
	})

	m = press(m, "o")
	m = press(m, "j")
	m = press(m, "t")

	if got := len(m.Tabs); got != 2 {
		t.Fatalf("len(Tabs) = %d, want 2", got)
	}
	if got := m.currentTab().Page.Path; got != "two.md" {
		t.Fatalf("current path = %q, want two.md", got)
	}
}

func TestSwitchTabsWithGtAndGT(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: []string{"a"}})
	m.Tabs = append(m.Tabs,
		Tab{Page: render.Page{Path: "b.md", Content: []string{"b"}}, History: nav.NewHistory(nav.HistoryEntry{Path: "b.md"})},
		Tab{Page: render.Page{Path: "c.md", Content: []string{"c"}}, History: nav.NewHistory(nav.HistoryEntry{Path: "c.md"})},
	)
	m.Tabs[0].ScrollY = 3

	m = pressSequence(m, "g", "t")
	if got := m.CurrentTab; got != 1 {
		t.Fatalf("CurrentTab after gt = %d, want 1", got)
	}
	if got := m.Tabs[0].ScrollY; got != 3 {
		t.Fatalf("first tab ScrollY after gt = %d, want 3", got)
	}

	m = pressSequence(m, "g", "T")
	if got := m.CurrentTab; got != 0 {
		t.Fatalf("CurrentTab after gT = %d, want 0", got)
	}

	m = pressSequence(m, "g", "T")
	if got := m.CurrentTab; got != 2 {
		t.Fatalf("CurrentTab after wrapped gT = %d, want 2", got)
	}
}

func TestCloseCurrentTab(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: []string{"a"}})
	m.Tabs = append(m.Tabs,
		Tab{Page: render.Page{Path: "b.md", Content: []string{"b"}}, History: nav.NewHistory(nav.HistoryEntry{Path: "b.md"})},
		Tab{Page: render.Page{Path: "c.md", Content: []string{"c"}}, History: nav.NewHistory(nav.HistoryEntry{Path: "c.md"})},
	)
	m.CurrentTab = 1

	m = press(m, "x")

	if got := len(m.Tabs); got != 2 {
		t.Fatalf("len(Tabs) = %d, want 2", got)
	}
	if got := m.CurrentTab; got != 1 {
		t.Fatalf("CurrentTab = %d, want 1", got)
	}
	if got := m.currentTab().Page.Path; got != "c.md" {
		t.Fatalf("current path = %q, want c.md", got)
	}
}

func TestCannotCloseLastTab(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: []string{"a"}})

	m = press(m, "x")

	if got := len(m.Tabs); got != 1 {
		t.Fatalf("len(Tabs) = %d, want 1", got)
	}
	if got := m.Status; got != "cannot close last tab" {
		t.Fatalf("Status = %q, want cannot close last tab", got)
	}
}

func TestRenderTabsShowsCurrentTab(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: []string{"a"}})
	m.Width = 80
	m.Tabs = append(m.Tabs,
		Tab{Page: render.Page{Path: "b.md", Content: []string{"b"}}, History: nav.NewHistory(nav.HistoryEntry{Path: "b.md"})},
	)
	m.CurrentTab = 1

	view := m.View()

	for _, text := range []string{"1 a.md", "2 b.md"} {
		if !strings.Contains(view, text) {
			t.Fatalf("View() does not show %q in tab list: %q", text, view)
		}
	}
	if !strings.Contains(view, "tab 2/2") {
		t.Fatalf("View() does not show tab status: %q", view)
	}
}

func TestRenderTabsTruncatesLongNames(t *testing.T) {
	m := testModel(render.Page{Path: "very-long-document-name.md", Content: []string{"a"}})
	m.Width = 20

	view := m.View()

	if !strings.Contains(view, "…") {
		t.Fatalf("View() does not truncate long tab name: %q", view)
	}
}

func TestWatchEventReloadsMatchingTabsAndKeepsScroll(t *testing.T) {
	loads := 0
	updated := render.Page{Path: "a.md", Content: []string{"updated", "line", "line", "line"}}
	m := NewModel(render.Page{Path: "a.md", Content: lines(4)}, func(path string) (render.Page, error) {
		loads++
		if path != "a.md" {
			t.Fatalf("path = %q, want a.md", path)
		}
		return updated, nil
	})
	m.Height = 4
	m.Tabs[0].ScrollY = 1
	m.Tabs = append(m.Tabs, Tab{
		Page:    render.Page{Path: "a.md", Content: lines(4)},
		History: nav.NewHistory(nav.HistoryEntry{Path: "a.md"}),
		ScrollY: 2,
	})
	m.WatchedPaths["a.md"] = true

	next, _ := m.Update(watch.Event{Path: "a.md"})
	m = next.(Model)

	if loads != 1 {
		t.Fatalf("loads = %d, want 1", loads)
	}
	for i, tab := range m.Tabs {
		if got := tab.Page.Content[0]; got != "updated" {
			t.Fatalf("tab %d content[0] = %q, want updated", i, got)
		}
	}
	if got := m.Tabs[0].ScrollY; got != 1 {
		t.Fatalf("tab 0 ScrollY = %d, want 1", got)
	}
	if got := m.Tabs[1].ScrollY; got != 2 {
		t.Fatalf("tab 1 ScrollY = %d, want 2", got)
	}
	if got := m.Status; got != "reloaded" {
		t.Fatalf("Status = %q, want reloaded", got)
	}
}

func TestWatchEventClampsScrollWhenReloadedFileShrinks(t *testing.T) {
	m := NewModel(render.Page{Path: "a.md", Content: lines(10)}, func(string) (render.Page, error) {
		return render.Page{Path: "a.md", Content: lines(2)}, nil
	})
	m.Height = 5
	m.Tabs[0].ScrollY = 7
	m.WatchedPaths["a.md"] = true

	next, _ := m.Update(watch.Event{Path: "a.md"})
	m = next.(Model)

	if got := m.currentTab().ScrollY; got != 0 {
		t.Fatalf("ScrollY = %d, want 0", got)
	}
}

func TestWatchEventReloadErrorKeepsPage(t *testing.T) {
	m := NewModel(render.Page{Path: "a.md", Content: []string{"old"}}, func(string) (render.Page, error) {
		return render.Page{}, errors.New("reload failed")
	})
	m.WatchedPaths["a.md"] = true

	next, _ := m.Update(watch.Event{Path: "a.md"})
	m = next.(Model)

	if got := m.currentTab().Page.Content[0]; got != "old" {
		t.Fatalf("content[0] = %q, want old", got)
	}
	if got := m.Status; got != "reload failed" {
		t.Fatalf("Status = %q, want reload failed", got)
	}
}

func TestWatchErrorUpdatesStatus(t *testing.T) {
	m := testModel(render.Page{Path: "a.md", Content: []string{"a"}})
	m.WatchedPaths["a.md"] = true

	next, _ := m.Update(watch.Error{Path: "a.md", Err: errors.New("watch failed")})
	m = next.(Model)

	if got := m.Status; got != "watch failed" {
		t.Fatalf("Status = %q, want watch failed", got)
	}
	if m.WatchedPaths["a.md"] {
		t.Fatal("WatchedPaths[a.md] is still true")
	}
}

func testModel(page render.Page) Model {
	m := NewModel(page, func(path string) (render.Page, error) {
		return render.Page{Path: path, Content: []string{path}}, nil
	})
	m.Height = 10
	return m
}

func press(m Model, s string) Model {
	return keyRunes(m, []rune(s))
}

func typeRunes(m Model, s string) Model {
	return keyRunes(m, []rune(s))
}

func pressSequence(m Model, keys ...string) Model {
	for _, key := range keys {
		m = press(m, key)
	}
	return m
}

func keyRunes(m Model, runes []rune) Model {
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: runes})
	return next.(Model)
}

func key(m Model, typ tea.KeyType) Model {
	next, _ := m.Update(tea.KeyMsg{Type: typ})
	return next.(Model)
}

func lines(n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = "line"
	}
	return out
}
