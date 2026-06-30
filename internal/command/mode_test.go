package command

import (
	"testing"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

type returnStub struct{}

func (returnStub) Init(_ *core.Context) tea.Cmd                           { return nil }
func (returnStub) Update(_ *core.Context, _ tea.Msg) (core.Mode, tea.Cmd) { return nil, nil }
func (returnStub) View(_ *core.Context) string                            { return "" }
func (returnStub) Help() help.KeyMap                                      { return nil }

func newCtx() *core.Context {
	tbl := table.New(table.WithFocused(true), table.WithStyles(core.TableStyles()))
	return &core.Context{
		Width:  120,
		Height: 30,
		Config: config.DefaultConfig(),
		Table:  tbl,
		Keys:   core.DefaultKeys(),
	}
}

func specialKey(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func TestCommandNew(t *testing.T) {
	c := New(returnStub{})
	if c == nil {
		t.Fatal("expected non-nil command")
	}
}

func TestCommandCancelReturnsToTarget(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	next, cmd := c.Update(ctx, specialKey(tea.KeyEscape))
	if cmd != nil {
		t.Fatal("expected nil cmd on cancel")
	}
	if _, ok := next.(returnStub); !ok {
		t.Fatalf("expected returnStub on cancel, got %T", next)
	}
}

func TestCommandEmptyLineReturnsToTarget(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("")
	next, cmd := c.Update(ctx, specialKey(tea.KeyEnter))
	if cmd != nil {
		t.Fatal("expected nil cmd on empty line")
	}
	if _, ok := next.(returnStub); !ok {
		t.Fatalf("expected returnStub on empty line, got %T", next)
	}
}

func TestCommandQuitTriggersQuit(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("quit")
	_, cmd := c.Update(ctx, specialKey(tea.KeyEnter))
	if cmd == nil {
		t.Fatal("expected quit cmd on quit")
	}
}

func TestCommandQAliasQuits(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("q")
	_, cmd := c.Update(ctx, specialKey(tea.KeyEnter))
	if cmd == nil {
		t.Fatal("expected quit cmd for :q")
	}
}

func TestCommandSortByName(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("sortby name")
	_, _ = c.Update(ctx, specialKey(tea.KeyEnter))
	if ctx.Config.SortBy != config.SortByName {
		t.Fatalf("expected SortByName, got %v", ctx.Config.SortBy)
	}
}

func TestCommandSortBySize(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("sortby size")
	_, _ = c.Update(ctx, specialKey(tea.KeyEnter))
	if ctx.Config.SortBy != config.SortBySize {
		t.Fatalf("expected SortBySize, got %v", ctx.Config.SortBy)
	}
}

func TestCommandSortAscending(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("sort ascending")
	_, _ = c.Update(ctx, specialKey(tea.KeyEnter))
	if ctx.Config.SortOrder != config.Ascending {
		t.Fatalf("expected Ascending, got %v", ctx.Config.SortOrder)
	}
}

func TestCommandSortDescending(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("sort descending")
	_, _ = c.Update(ctx, specialKey(tea.KeyEnter))
	if ctx.Config.SortOrder != config.Descending {
		t.Fatalf("expected Descending, got %v", ctx.Config.SortOrder)
	}
}

func TestCommandSortAscAlias(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("sort asc")
	_, _ = c.Update(ctx, specialKey(tea.KeyEnter))
	if ctx.Config.SortOrder != config.Ascending {
		t.Fatalf("expected Ascending with 'asc' alias, got %v", ctx.Config.SortOrder)
	}
}

func TestCommandSortDescAlias(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("sort desc")
	_, _ = c.Update(ctx, specialKey(tea.KeyEnter))
	if ctx.Config.SortOrder != config.Descending {
		t.Fatalf("expected Descending with 'desc' alias, got %v", ctx.Config.SortOrder)
	}
}

func TestCommandGroupFilesFirst(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("group files")
	_, _ = c.Update(ctx, specialKey(tea.KeyEnter))
	if ctx.Config.Grouping != config.FilesFirst {
		t.Fatalf("expected FilesFirst, got %v", ctx.Config.Grouping)
	}
}

func TestCommandGroupMixed(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("group mixed")
	_, _ = c.Update(ctx, specialKey(tea.KeyEnter))
	if ctx.Config.Grouping != config.Mixed {
		t.Fatalf("expected Mixed, got %v", ctx.Config.Grouping)
	}
}

func TestCommandToggleIcons(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	prev := ctx.Config.ShowIcons
	c.input.SetValue("toggle icons")
	_, _ = c.Update(ctx, specialKey(tea.KeyEnter))
	if ctx.Config.ShowIcons == prev {
		t.Fatal("expected icons toggle")
	}
}

func TestCommandToggleHidden(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	prev := ctx.Config.ShowHidden
	c.input.SetValue("toggle hidden")
	_, _ = c.Update(ctx, specialKey(tea.KeyEnter))
	if ctx.Config.ShowHidden == prev {
		t.Fatal("expected hidden toggle")
	}
}

func TestCommandToggleLive(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	ctx.Config.LiveFilter = false
	c.input.SetValue("toggle live")
	_, _ = c.Update(ctx, specialKey(tea.KeyEnter))
	if !ctx.Config.LiveFilter {
		t.Fatal("expected live filter toggle to true")
	}
}

func TestCommandUnknownVerbReturnsToTarget(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("nonexistent")
	next, cmd := c.Update(ctx, specialKey(tea.KeyEnter))
	if cmd != nil {
		t.Fatal("expected nil cmd on unknown verb")
	}
	if _, ok := next.(returnStub); !ok {
		t.Fatalf("expected returnStub on unknown verb, got %T", next)
	}
}

func TestCommandVerbWithoutArgReturnsToTarget(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	c.input.SetValue("sort")
	next, cmd := c.Update(ctx, specialKey(tea.KeyEnter))
	if cmd != nil {
		t.Fatal("expected nil cmd on verb without arg")
	}
	if _, ok := next.(returnStub); !ok {
		t.Fatalf("expected returnStub on verb without arg, got %T", next)
	}
}

func TestCommandViewNonEmpty(t *testing.T) {
	c := New(returnStub{})
	ctx := newCtx()
	_ = c.Init(ctx)
	v := c.View(ctx)
	if v == "" {
		t.Fatal("expected non-empty view")
	}
}

func TestCommandHelpReturnsBindings(t *testing.T) {
	c := New(returnStub{})
	bindings := c.Help()
	if bindings == nil {
		t.Fatal("expected non-nil help")
	}
	if len(bindings.ShortHelp()) == 0 {
		t.Fatal("expected non-empty help")
	}
}
