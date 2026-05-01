package app

import (
	"testing"

	"github.com/YamasouA/mdview/internal/render"
)

func TestNewModelWithPagesCreatesInitialTabs(t *testing.T) {
	m := NewModelWithPages([]render.Page{
		{Path: "README.md", Content: []string{"readme"}},
		{Path: "README.en.md", Content: []string{"english"}},
	}, func(path string) (render.Page, error) {
		return render.Page{Path: path}, nil
	})

	if got := len(m.Tabs); got != 2 {
		t.Fatalf("len(Tabs) = %d, want 2", got)
	}
	if got := m.CurrentTab; got != 0 {
		t.Fatalf("CurrentTab = %d, want 0", got)
	}
	if got := m.Tabs[0].Page.Path; got != "README.md" {
		t.Fatalf("Tabs[0].Path = %q, want README.md", got)
	}
	if got := m.Tabs[1].Page.Path; got != "README.en.md" {
		t.Fatalf("Tabs[1].Path = %q, want README.en.md", got)
	}
}
