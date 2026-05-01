package render

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

func TestExtractLinks(t *testing.T) {
	raw := "# Title\nSee [Setup](docs/setup.md) and [Site](https://example.com).\n"

	links := ExtractLinks(raw)

	if len(links) != 2 {
		t.Fatalf("len(links) = %d, want 2", len(links))
	}
	if links[0].Text != "Setup" || links[0].Target != "docs/setup.md" || links[0].Line != 1 {
		t.Fatalf("links[0] = %+v", links[0])
	}
	if links[1].Text != "Site" || links[1].Target != "https://example.com" || links[1].Line != 1 {
		t.Fatalf("links[1] = %+v", links[1])
	}
}

func TestExtractLinksIgnoresFencedCodeBlocks(t *testing.T) {
	raw := "See [Real](real.md).\n\n```md\n[Example](missing.md)\n```\n"

	links := ExtractLinks(raw)

	if len(links) != 1 {
		t.Fatalf("len(links) = %d, want 1", len(links))
	}
	if links[0].Text != "Real" || links[0].Target != "real.md" {
		t.Fatalf("links[0] = %+v", links[0])
	}
}

func TestExtractHeadings(t *testing.T) {
	raw := "# Title\n\nbody\n### Details ###\n####### not heading\n#not heading\n"

	headings := ExtractHeadings(raw)

	if len(headings) != 2 {
		t.Fatalf("len(headings) = %d, want 2", len(headings))
	}
	if headings[0].Text != "Title" || headings[0].Level != 1 || headings[0].Line != 0 {
		t.Fatalf("headings[0] = %+v", headings[0])
	}
	if headings[1].Text != "Details" || headings[1].Level != 3 || headings[1].Line != 3 {
		t.Fatalf("headings[1] = %+v", headings[1])
	}
}

func TestExtractHeadingsIgnoresFencedCodeBlocks(t *testing.T) {
	raw := "# Real\n\n```md\n## Example\n```\n"

	headings := ExtractHeadings(raw)

	if len(headings) != 1 {
		t.Fatalf("len(headings) = %d, want 1", len(headings))
	}
	if headings[0].Text != "Real" || headings[0].Level != 1 {
		t.Fatalf("headings[0] = %+v", headings[0])
	}
}

func TestRenderMarkdownUsesGlamourDarkHeadingStyle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	if err := os.WriteFile(path, []byte("# h1\n\n## 操作\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	page, err := RenderMarkdown(path)
	if err != nil {
		t.Fatal(err)
	}

	plain := stripANSI(strings.Join(page.Content, "\n"))
	if strings.Contains(plain, "# h1") {
		t.Fatalf("h1 marker should be removed by glamour dark style: %q", plain)
	}
	if !strings.Contains(plain, "## 操作") {
		t.Fatalf("h2 marker should follow glamour dark style: %q", plain)
	}
	for _, heading := range []string{"h1", "操作"} {
		if !strings.Contains(plain, heading) {
			t.Fatalf("heading %q was not rendered in %q", heading, plain)
		}
	}
}

func TestResolveTargetRelativeToCurrentFile(t *testing.T) {
	got := ResolveTarget(filepath.Join("docs", "guide", "intro.md"), "../setup.md")
	want := filepath.Join("docs", "setup.md")
	if got != want {
		t.Fatalf("ResolveTarget() = %q, want %q", got, want)
	}
}

func TestRenderMarkdown(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	if err := os.WriteFile(path, []byte("# Hello\n\n[Next](next.md)\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	page, err := RenderMarkdown(path)
	if err != nil {
		t.Fatal(err)
	}
	if page.Path != path {
		t.Fatalf("Path = %q, want %q", page.Path, path)
	}
	if len(page.Content) == 0 {
		t.Fatal("Content is empty")
	}
	if len(page.Links) != 1 {
		t.Fatalf("len(Links) = %d, want 1", len(page.Links))
	}
	if len(page.Headings) != 1 {
		t.Fatalf("len(Headings) = %d, want 1", len(page.Headings))
	}
	if got := page.Headings[0].Text; got != "Hello" {
		t.Fatalf("Headings[0].Text = %q, want Hello", got)
	}
}

func stripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}
