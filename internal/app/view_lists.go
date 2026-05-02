package app

import (
	"fmt"
	"strings"

	"github.com/YamasouA/mdview/internal/render"
)

func (m Model) markEnterLinkLine(lines []string, start int) {
	links := m.currentTab().Page.Links
	if len(links) == 0 {
		return
	}
	line := links[0].Line
	if line < start || line >= start+len(lines) {
		return
	}
	lines[line-start] = "enter> " + lines[line-start]
}

func matchesOnLine(matches []SearchMatch, line int) []SearchMatch {
	var out []SearchMatch
	for _, match := range matches {
		if match.Line == line {
			out = append(out, match)
		}
	}
	return out
}

func indexOnLine(matches []SearchMatch, index int) int {
	if index < 0 || index >= len(matches) {
		return -1
	}
	line := matches[index].Line
	active := 0
	for i := 0; i < index; i++ {
		if matches[i].Line == line {
			active++
		}
	}
	return active
}

func (m Model) renderLinks(height int) string {
	links := m.currentTab().Page.Links
	lines := make([]string, 0, len(links))
	for i, link := range links {
		prefix := "  "
		if i == m.LinkIndex {
			prefix = "> "
		}
		resolved := render.ResolveTarget(m.currentTab().Page.Path, link.Target)
		lines = append(lines, fmt.Sprintf("%s%d. line %d | %s -> %s | resolved: %s", prefix, i+1, link.Line+1, link.Text, link.Target, resolved))
	}
	if len(lines) > height {
		lines = lines[:height]
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderHeadings(height int) string {
	headings := m.currentTab().Page.Headings
	lines := make([]string, 0, len(headings))
	for i, heading := range headings {
		prefix := "  "
		if i == m.HeadingIndex {
			prefix = "> "
		}
		lines = append(lines, fmt.Sprintf("%s%s %s", prefix, strings.Repeat("#", heading.Level), heading.Text))
	}
	if len(lines) > height {
		lines = lines[:height]
	}
	for len(lines) < height {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}
