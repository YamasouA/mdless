package app

import (
	"strings"
	"unicode/utf8"
)

const (
	ansiEscapeStart      = '\x1b'
	highlightStart       = "\x1b[7m"
	activeHighlightStart = "\x1b[1;30;43m"
	highlightEnd         = "\x1b[0m"
)

func plainText(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if next, ok := skipANSI(s, i); ok {
			i = next
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		b.WriteRune(r)
		i += size
	}
	return b.String()
}

func containsFoldPlain(line, query string) bool {
	return len(findPlainMatches(0, line, query)) > 0
}

func findPlainMatches(lineNo int, line, query string) []SearchMatch {
	if query == "" {
		return nil
	}

	plain, _ := plainTextWithRawOffsets(line)
	lowerPlain := strings.ToLower(plain)
	lowerQuery := strings.ToLower(query)
	if !strings.Contains(lowerPlain, lowerQuery) {
		return nil
	}

	var matches []SearchMatch
	searchPos := 0
	for {
		matchStart := strings.Index(lowerPlain[searchPos:], lowerQuery)
		if matchStart < 0 {
			break
		}
		plainStart := searchPos + matchStart
		plainEnd := plainStart + len(lowerQuery)
		matches = append(matches, SearchMatch{
			Line:  lineNo,
			Start: plainStart,
			End:   plainEnd,
		})
		searchPos = plainEnd
	}
	return matches
}

func highlightMatches(line string, matches []SearchMatch, activeIndex int) string {
	if len(matches) == 0 {
		return line
	}

	_, plainToRaw := plainTextWithRawOffsets(line)

	var out strings.Builder
	rawPos := 0
	for i, match := range matches {
		rawStart := plainToRaw[match.Start]
		rawEnd := plainToRaw[match.End]

		out.WriteString(line[rawPos:rawStart])
		if i == activeIndex {
			out.WriteString(activeHighlightStart)
		} else {
			out.WriteString(highlightStart)
		}
		out.WriteString(plainText(line[rawStart:rawEnd]))
		out.WriteString(highlightEnd)

		rawPos = rawEnd
	}
	out.WriteString(line[rawPos:])
	return out.String()
}

func plainTextWithRawOffsets(s string) (string, []int) {
	var plain strings.Builder
	offsets := []int{0}
	for i := 0; i < len(s); {
		if next, ok := skipANSI(s, i); ok {
			i = next
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		plain.WriteRune(r)
		for range size {
			offsets = append(offsets, i+size)
		}
		i += size
	}
	return plain.String(), offsets
}

func skipANSI(s string, i int) (int, bool) {
	if i >= len(s) || s[i] != ansiEscapeStart || i+1 >= len(s) || s[i+1] != '[' {
		return i, false
	}
	for j := i + 2; j < len(s); j++ {
		if s[j] >= '@' && s[j] <= '~' {
			return j + 1, true
		}
	}
	return len(s), true
}
