package nav

type History struct {
	Entries []HistoryEntry
	Index   int
}

type HistoryEntry struct {
	Path    string
	ScrollY int
}

func NewHistory(entry HistoryEntry) History {
	return History{
		Entries: []HistoryEntry{entry},
		Index:   0,
	}
}

func (h History) Push(entry HistoryEntry) History {
	if len(h.Entries) == 0 {
		return NewHistory(entry)
	}

	entries := append([]HistoryEntry(nil), h.Entries[:h.Index+1]...)
	entries = append(entries, entry)

	return History{
		Entries: entries,
		Index:   len(entries) - 1,
	}
}

func (h History) UpdateCurrent(scrollY int) History {
	if len(h.Entries) == 0 || h.Index < 0 || h.Index >= len(h.Entries) {
		return h
	}

	entries := append([]HistoryEntry(nil), h.Entries...)
	entries[h.Index].ScrollY = scrollY
	h.Entries = entries
	return h
}

func (h History) Back() (History, HistoryEntry, bool) {
	if len(h.Entries) == 0 || h.Index <= 0 {
		return h, HistoryEntry{}, false
	}

	h.Index--
	return h, h.Entries[h.Index], true
}

func (h History) Forward() (History, HistoryEntry, bool) {
	if len(h.Entries) == 0 || h.Index >= len(h.Entries)-1 {
		return h, HistoryEntry{}, false
	}

	h.Index++
	return h, h.Entries[h.Index], true
}
