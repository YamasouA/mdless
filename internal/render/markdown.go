package render

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/YamasouA/mdless/internal/nav"
	"github.com/charmbracelet/glamour"
)

type Page struct {
	Path    string
	Content []string
	Links   []nav.Link
}

var inlineLinkPattern = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

func RenderMarkdown(path string) (Page, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Page{}, err
	}

	out, err := glamour.Render(string(raw), "dark")
	if err != nil {
		return Page{}, err
	}

	return Page{
		Path:    path,
		Content: strings.Split(strings.TrimRight(out, "\n"), "\n"),
		Links:   ExtractLinks(string(raw)),
	}, nil
}

func ExtractLinks(raw string) []nav.Link {
	var links []nav.Link
	lines := strings.Split(raw, "\n")
	inFence := false
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}

		matches := inlineLinkPattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			links = append(links, nav.Link{
				Text:   match[1],
				Target: match[2],
				Line:   i,
			})
		}
	}
	return links
}

func ResolveTarget(currentPath, target string) string {
	if filepath.IsAbs(target) || strings.Contains(target, "://") {
		return target
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentPath), target))
}
