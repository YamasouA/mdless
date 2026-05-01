package app

import (
	"strings"
	"testing"
)

func TestPlainTextStripsANSI(t *testing.T) {
	got := plainText("\x1b[31mTarget\x1b[0m")
	if got != "Target" {
		t.Fatalf("plainText() = %q, want Target", got)
	}
}

func TestContainsFoldPlainIgnoresANSI(t *testing.T) {
	if !containsFoldPlain("\x1b[31mTarget\x1b[0m", "target") {
		t.Fatal("containsFoldPlain() returned false")
	}
}

func TestHighlightMatchesPreservesTextAndAddsHighlight(t *testing.T) {
	line := "\x1b[31mTarget\x1b[0m target"
	matches := findPlainMatches(0, line, "target")
	got := highlightMatches(line, matches, -1)
	plain := plainText(got)

	if plain != "Target target" {
		t.Fatalf("plain text = %q, want Target target", plain)
	}
	if count := strings.Count(got, highlightStart); count != 2 {
		t.Fatalf("highlight count = %d, want 2 in %q", count, got)
	}
}

func TestHighlightMatchesMarksActiveMatchDifferently(t *testing.T) {
	line := "target target"
	matches := findPlainMatches(0, line, "target")

	got := highlightMatches(line, matches, 1)

	if count := strings.Count(got, highlightStart); count != 1 {
		t.Fatalf("normal highlight count = %d, want 1 in %q", count, got)
	}
	if count := strings.Count(got, activeHighlightStart); count != 1 {
		t.Fatalf("active highlight count = %d, want 1 in %q", count, got)
	}
}

func TestHighlightMatchesDoesNotLoseHighlightAcrossANSIReset(t *testing.T) {
	line := "\x1b[31mtar\x1b[0mget"
	matches := findPlainMatches(0, line, "target")

	got := highlightMatches(line, matches, 0)

	if !strings.Contains(got, activeHighlightStart+"target"+highlightEnd) {
		t.Fatalf("active highlight does not cover full match across ANSI reset: %q", got)
	}
	if plain := plainText(got); plain != "target" {
		t.Fatalf("plain text = %q, want target", plain)
	}
}
