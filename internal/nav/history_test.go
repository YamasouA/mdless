package nav

import "testing"

func TestHistoryBackForwardEmpty(t *testing.T) {
	var h History

	if _, _, ok := h.Back(); ok {
		t.Fatal("Back() on empty history returned ok")
	}
	if _, _, ok := h.Forward(); ok {
		t.Fatal("Forward() on empty history returned ok")
	}
}

func TestHistoryPushBackForward(t *testing.T) {
	h := NewHistory(HistoryEntry{Path: "a.md", ScrollY: 1})
	h = h.Push(HistoryEntry{Path: "b.md", ScrollY: 2})

	var entry HistoryEntry
	var ok bool
	h, entry, ok = h.Back()
	if !ok {
		t.Fatal("Back() returned false")
	}
	if entry.Path != "a.md" || entry.ScrollY != 1 {
		t.Fatalf("Back() entry = %+v", entry)
	}

	h, entry, ok = h.Forward()
	if !ok {
		t.Fatal("Forward() returned false")
	}
	if entry.Path != "b.md" || entry.ScrollY != 2 {
		t.Fatalf("Forward() entry = %+v", entry)
	}
}

func TestHistoryPushAfterBackDropsForwardEntries(t *testing.T) {
	h := NewHistory(HistoryEntry{Path: "a.md"})
	h = h.Push(HistoryEntry{Path: "b.md"})
	h = h.Push(HistoryEntry{Path: "c.md"})

	var ok bool
	h, _, ok = h.Back()
	if !ok {
		t.Fatal("Back() returned false")
	}

	h = h.Push(HistoryEntry{Path: "d.md"})

	if len(h.Entries) != 3 {
		t.Fatalf("len(Entries) = %d, want 3", len(h.Entries))
	}
	if h.Entries[2].Path != "d.md" {
		t.Fatalf("last entry = %+v", h.Entries[2])
	}
	if _, _, ok := h.Forward(); ok {
		t.Fatal("Forward() after branch returned ok")
	}
}

func TestHistoryUpdateCurrent(t *testing.T) {
	h := NewHistory(HistoryEntry{Path: "a.md"})
	h = h.Push(HistoryEntry{Path: "b.md"})

	h = h.UpdateCurrent(12)

	if got := h.Entries[h.Index].ScrollY; got != 12 {
		t.Fatalf("current ScrollY = %d, want 12", got)
	}
	if got := h.Entries[0].ScrollY; got != 0 {
		t.Fatalf("first ScrollY = %d, want 0", got)
	}
}
