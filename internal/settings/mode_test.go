package settings

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	stexmodel "github.com/SolracHQ/stex/internal/model"
	"github.com/SolracHQ/stex/internal/testutil"
	"github.com/SolracHQ/stex/internal/styles"

	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

// buildTestTree creates a small tree in a temp dir, scans it, and returns the root.
func buildTestTree(t *testing.T) *stexmodel.Dir {
	t.Helper()
	tmp := t.TempDir()
	for _, name := range []string{"alpha.txt", "beta.txt", "gamma.txt"} {
		if err := os.WriteFile(filepath.Join(tmp, name), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	state := &stexmodel.ScanState{}
	stexmodel.BuildTree(tmp, state)
	return state.Result
}

func newCtx(t *testing.T) *core.Context {
	t.Helper()
	root := buildTestTree(t)
	tbl := table.New(table.WithFocused(true), table.WithStyles(styles.TableDefault()))
	ctx := &core.Context{
		Width:   120,
		Height:  30,
		Root:    root,
		Current: root,
		Config:  config.DefaultConfig(),
		Table:   tbl,
		Keys:    core.DefaultKeys(),
	}
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	return ctx
}

func specialKey(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func TestSettingsNew(t *testing.T) {
	s := New(testutil.StubMode{})
	if s == nil {
		t.Fatal("expected non-nil settings")
	}
}

func TestSettingsInitTakesSnapshot(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	prev := ctx.Config.SortBy
	ctx.Config.SortBy = config.SortBySize
	_ = s.Init(ctx)
	_, _ = s.Update(ctx, specialKey(tea.KeyEscape))
	if ctx.Config.SortBy != prev {
		t.Fatal("expected config to be reverted on esc")
	}
}

func TestSettingsDownMovesCursor(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	_, _ = s.Update(ctx, runeKey('j'))
	if s.cursor != 1 {
		t.Fatalf("expected cursor at 1, got %d", s.cursor)
	}
}

func TestSettingsDownClampsAtLast(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	for i := 0; i < rowCount+2; i++ {
		_, _ = s.Update(ctx, runeKey('j'))
	}
	if s.cursor != rowCount-1 {
		t.Fatalf("expected cursor clamped at %d, got %d", rowCount-1, s.cursor)
	}
}

func TestSettingsUpMovesCursor(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	s.cursor = 2
	_, _ = s.Update(ctx, runeKey('k'))
	if s.cursor != 1 {
		t.Fatalf("expected cursor at 1, got %d", s.cursor)
	}
}

func TestSettingsUpClampsAtFirst(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	_, _ = s.Update(ctx, runeKey('k'))
	if s.cursor != 0 {
		t.Fatalf("expected cursor at 0, got %d", s.cursor)
	}
}

func TestSettingsTabTogglesSort(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	prev := ctx.Config.SortBy
	_, _ = s.Update(ctx, specialKey(tea.KeyTab))
	if ctx.Config.SortBy == prev {
		t.Fatal("expected sort to toggle on tab")
	}
}

func TestSettingsTabTogglesOrder(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	s.cursor = rowOrder
	prev := ctx.Config.SortOrder
	_, _ = s.Update(ctx, specialKey(tea.KeyTab))
	if ctx.Config.SortOrder == prev {
		t.Fatal("expected order to toggle on tab")
	}
}

func TestSettingsTabOnGroupOpensPicker(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	s.cursor = rowGroup
	next, _ := s.Update(ctx, specialKey(tea.KeyTab))
	if next == nil {
		t.Fatal("expected grouping picker on tab for group row, got nil")
	}
}

func TestSettingsTabTogglesIcons(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	s.cursor = rowIcons
	prev := ctx.Config.ShowIcons
	_, _ = s.Update(ctx, specialKey(tea.KeyTab))
	if ctx.Config.ShowIcons == prev {
		t.Fatal("expected icons to toggle on tab")
	}
}

func TestSettingsTabTogglesHidden(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	s.cursor = rowHidden
	prev := ctx.Config.ShowHidden
	_, _ = s.Update(ctx, specialKey(tea.KeyTab))
	if ctx.Config.ShowHidden == prev {
		t.Fatal("expected hidden to toggle on tab")
	}
}

func TestSettingsTabTogglesLiveFilter(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	s.cursor = rowLiveFilter
	prev := ctx.Config.LiveFilter
	_, _ = s.Update(ctx, specialKey(tea.KeyTab))
	if ctx.Config.LiveFilter == prev {
		t.Fatal("expected live filter to toggle on tab")
	}
}

func TestSettingsConfirmReturnsToTarget(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	next, _ := s.Update(ctx, specialKey(tea.KeyEnter))
	if _, ok := next.(testutil.StubMode); !ok {
		t.Fatalf("expected testutil.StubMode on enter, got %T", next)
	}
}

func TestSettingsEscRevertsConfig(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	orig := ctx.Config.SortBy
	ctx.Config.SortBy = config.SortByName
	_, _ = s.Update(ctx, specialKey(tea.KeyEscape))
	if ctx.Config.SortBy != orig {
		t.Fatal("expected config to revert to snapshot on esc")
	}
}

func TestSettingsEscReturnsToTarget(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	next, _ := s.Update(ctx, specialKey(tea.KeyEscape))
	if _, ok := next.(testutil.StubMode); !ok {
		t.Fatalf("expected testutil.StubMode on esc, got %T", next)
	}
}

func TestSettingsSaveReturnsToTarget(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	next, _ := s.Update(ctx, runeKey('S'))
	if _, ok := next.(testutil.StubMode); !ok {
		t.Fatalf("expected testutil.StubMode on S, got %T", next)
	}
}

func TestSettingsResetToDefaults(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	ctx.Config.SortBy = config.SortByName
	_, _ = s.Update(ctx, runeKey('r'))
	if ctx.Config.SortBy != config.SortBySize {
		t.Fatal("expected sort to be reset to size (default)")
	}
}

func TestSettingsViewNonEmpty(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	v := s.Overlay(ctx)
	if v == "" {
		t.Fatal("expected non-empty view")
	}
}

func TestSettingsHelpReturnsBindings(t *testing.T) {
	s := New(testutil.StubMode{})
	bindings := s.Help()
	if bindings == nil {
		t.Fatal("expected non-nil help")
	}
	if len(bindings.ShortHelp()) == 0 {
		t.Fatal("expected non-empty short help")
	}
}

func TestSettingsEnterKeepsChanges(t *testing.T) {
	ctx := newCtx(t)
	s := New(testutil.StubMode{})
	_ = s.Init(ctx)
	ctx.Config.ShowIcons = true
	_, _ = s.Update(ctx, specialKey(tea.KeyEnter))
	if !ctx.Config.ShowIcons {
		t.Fatal("expected ShowIcons to stay true after enter")
	}
}

func runeKey(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)})
}
