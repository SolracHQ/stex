package choose

import (
	"strings"
	"testing"

	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/testutil"

	tea "charm.land/bubbletea/v2"
)

func specialKey(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func TestChooseNew(t *testing.T) {
	c := New("Test?", []Option{
		{Label: "Yes"},
		{Label: "No"},
	}, nil)
	if c.cursor != 0 {
		t.Fatalf("expected cursor at 0, got %d", c.cursor)
	}
	if len(c.options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(c.options))
	}
}

func TestChooseDownMovesCursor(t *testing.T) {
	c := New("Test?", []Option{
		{Label: "A"},
		{Label: "B"},
	}, nil)
	_, _ = c.Update(&core.Context{}, runeKey('j'))
	if c.cursor != 1 {
		t.Fatalf("expected cursor at 1, got %d", c.cursor)
	}
}

func TestChooseDownClampsAtLast(t *testing.T) {
	c := New("Test?", []Option{
		{Label: "A"},
	}, nil)
	_, _ = c.Update(&core.Context{}, runeKey('j'))
	if c.cursor != 0 {
		t.Fatalf("expected cursor clamped at 0, got %d", c.cursor)
	}
}

func TestChooseUpClampsAtFirst(t *testing.T) {
	c := New("Test?", []Option{
		{Label: "A"},
	}, nil)
	_, _ = c.Update(&core.Context{}, runeKey('k'))
	if c.cursor != 0 {
		t.Fatalf("expected cursor clamped at 0, got %d", c.cursor)
	}
}

func TestChooseConfirmInvokesAction(t *testing.T) {
	called := false
	c := New("Test?", []Option{
		{
			Label: "Yes",
			Action: func(*core.Context) (core.Mode, tea.Cmd) {
				called = true
				return &testutil.StubMode{}, nil
			},
		},
	}, nil)
	_, _ = c.Update(&core.Context{}, specialKey(tea.KeyEnter))
	if !called {
		t.Fatal("expected action to be called on confirm")
	}
}

func TestChooseViewNonEmpty(t *testing.T) {
	c := New("Test?", []Option{
		{Label: "Yes"},
		{Label: "No"},
	}, nil)
	ctx := &core.Context{Width: 80, Height: 24}
	v := c.Overlay(ctx)
	if v == "" {
		t.Fatal("expected non-empty view")
	}
}

func TestChooseViewContainsOptions(t *testing.T) {
	c := New("Title text", []Option{
		{Label: "AlphaOption"},
		{Label: "BetaOption"},
	}, nil)
	ctx := &core.Context{Width: 80, Height: 24}
	v := c.Overlay(ctx)
	if !strings.Contains(v, "AlphaOption") {
		t.Error("expected view to contain AlphaOption")
	}
	if !strings.Contains(v, "BetaOption") {
		t.Error("expected view to contain BetaOption")
	}
}

func TestChooseHelpReturnsBindings(t *testing.T) {
	c := New("Test?", []Option{{Label: "X"}}, nil)
	bindings := c.Help()
	if bindings == nil {
		t.Fatal("expected non-nil help")
	}
	if len(bindings.ShortHelp()) == 0 {
		t.Fatal("expected non-empty help")
	}
}

func TestChooseCancelReturnsToBackTo(t *testing.T) {
	sentinel := &testutil.StubMode{}
	c := New("Test?", []Option{{Label: "X"}}, sentinel)
	next, _ := c.Update(&core.Context{}, specialKey(tea.KeyEscape))
	if next != sentinel {
		t.Fatalf("expected sentinel on cancel, got %T", next)
	}
}

func runeKey(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)})
}
