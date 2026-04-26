package app

import (
	"strings"
	"unicode/utf8"
)

const (
	ansiEscapeStart = '\x1b'
	highlightStart  = "\x1b[7m"
	highlightEnd    = "\x1b[0m"
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
	if query == "" {
		return false
	}
	return strings.Contains(strings.ToLower(plainText(line)), strings.ToLower(query))
}

func highlightMatches(line, query string) string {
	if query == "" {
		return line
	}

	plain, plainToRaw := plainTextWithRawOffsets(line)
	lowerPlain := strings.ToLower(plain)
	lowerQuery := strings.ToLower(query)
	if !strings.Contains(lowerPlain, lowerQuery) {
		return line
	}

	var out strings.Builder
	rawPos := 0
	searchPos := 0
	for {
		matchStart := strings.Index(lowerPlain[searchPos:], lowerQuery)
		if matchStart < 0 {
			break
		}
		plainStart := searchPos + matchStart
		plainEnd := plainStart + len(lowerQuery)
		rawStart := plainToRaw[plainStart]
		rawEnd := plainToRaw[plainEnd]

		out.WriteString(line[rawPos:rawStart])
		out.WriteString(highlightStart)
		out.WriteString(line[rawStart:rawEnd])
		out.WriteString(highlightEnd)

		rawPos = rawEnd
		searchPos = plainEnd
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
