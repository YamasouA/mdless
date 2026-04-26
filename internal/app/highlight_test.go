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
	got := highlightMatches("\x1b[31mTarget\x1b[0m target", "target")
	plain := plainText(got)

	if plain != "Target target" {
		t.Fatalf("plain text = %q, want Target target", plain)
	}
	if count := strings.Count(got, highlightStart); count != 2 {
		t.Fatalf("highlight count = %d, want 2 in %q", count, got)
	}
}
