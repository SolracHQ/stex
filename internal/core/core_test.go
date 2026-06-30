package core

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/key"
)

func TestDefaultKeysHasQuitAndHelp(t *testing.T) {
	keys := DefaultKeys()
	if len(keys.Quit.Keys()) == 0 {
		t.Fatal("expected quit binding to have at least one key")
	}
	if len(keys.Help.Keys()) == 0 {
		t.Fatal("expected help binding to have at least one key")
	}
}

func TestFlatKeyMapShortHelp(t *testing.T) {
	a := key.NewBinding(key.WithKeys("a"))
	b := key.NewBinding(key.WithKeys("b"))
	km := FlatKeyMap{a, b}
	short := km.ShortHelp()
	if len(short) != 2 {
		t.Fatalf("expected 2 bindings in short help, got %d", len(short))
	}
}

func TestFlatKeyMapFullHelp(t *testing.T) {
	a := key.NewBinding(key.WithKeys("a"))
	b := key.NewBinding(key.WithKeys("b"))
	km := FlatKeyMap{a, b}
	full := km.FullHelp()
	if len(full) != 1 {
		t.Fatalf("expected 1 row in full help, got %d", len(full))
	}
	if len(full[0]) != 2 {
		t.Fatalf("expected 2 bindings in the row, got %d", len(full[0]))
	}
}

func TestBlankEmptyDims(t *testing.T) {
	if Blank(0, 0) != "" {
		t.Fatal("expected empty for zero dims")
	}
	if Blank(10, 0) != "" {
		t.Fatal("expected empty for zero height")
	}
	if Blank(0, 5) != "" {
		t.Fatal("expected empty for zero width")
	}
}

func TestBlankRendersFullSize(t *testing.T) {
	out := Blank(4, 3)
	lines := strings.Split(out, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	for i, line := range lines {
		if len(line) != 4 {
			t.Fatalf("line %d: expected 4 chars, got %d", i, len(line))
		}
	}
}

func TestRenderBaseEmptyContext(t *testing.T) {
	ctx := &Context{}
	if got := RenderBase(ctx, nil); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}
