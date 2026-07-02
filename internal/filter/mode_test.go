package filter

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	stexmodel "github.com/SolracHQ/stex/internal/model"
	"github.com/SolracHQ/stex/internal/styles"
	"github.com/SolracHQ/stex/internal/testutil"

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

// newCtx returns a Context wired with the test tree and the default config.
func newCtx(t *testing.T) *core.Context {
	t.Helper()
	root := buildTestTree(t)
	tbl := table.New(table.WithFocused(true), table.WithStyles(styles.TableDefault()))
	cfg := config.DefaultConfig()
	cfg.LiveFilter = true
	return &core.Context{
		Width:   120,
		Height:  30,
		Root:    root,
		Current: root,
		Config:  cfg,
		Table:   tbl,
		Keys:    core.DefaultKeys(),
	}
}

var sentinel core.Mode = testutil.StubMode{}

func runeKey(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)})
}

func specialKey(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func TestFilterNew(t *testing.T) {
	f := New(sentinel)
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestFilterInitFocuses(t *testing.T) {
	f := New(sentinel)
	cmd := f.Init(&core.Context{})
	if cmd == nil {
		t.Fatal("expected init to return focus command")
	}
}

func TestFilterCancelClearsFilter(t *testing.T) {
	ctx := newCtx(t)
	re, _ := regexp.Compile("a")
	ctx.Config.Filter = re

	f := New(sentinel)
	_ = f.Init(ctx)
	next, _ := f.Update(ctx, specialKey(tea.KeyEscape))

	if ctx.Config.Filter != nil {
		t.Fatal("expected filter to be cleared on cancel")
	}
	if next != sentinel {
		t.Fatalf("expected returnTo on cancel, got %T", next)
	}
}

func TestFilterConfirmCommitsFilter(t *testing.T) {
	ctx := newCtx(t)
	ctx.Config.LiveFilter = false

	f := New(sentinel)
	_ = f.Init(ctx)
	f.input.SetValue("alpha")
	next, _ := f.Update(ctx, specialKey(tea.KeyEnter))

	if ctx.Config.Filter == nil {
		t.Fatal("expected filter to be committed on enter")
	}
	if !ctx.Config.Filter.MatchString("alpha.txt") {
		t.Fatal("expected filter to match alpha.txt")
	}
	if next != sentinel {
		t.Fatalf("expected returnTo on confirm, got %T", next)
	}
}

func TestFilterLiveModeAppliesOnKey(t *testing.T) {
	ctx := newCtx(t)
	ctx.Config.LiveFilter = true

	f := New(sentinel)
	_ = f.Init(ctx)
	_, _ = f.Update(ctx, runeKey('a'))

	if ctx.Config.Filter == nil {
		t.Fatal("expected live filter to commit on keystroke")
	}
}

func TestFilterManualModeDoesNotApplyOnKey(t *testing.T) {
	ctx := newCtx(t)
	ctx.Config.LiveFilter = false

	f := New(sentinel)
	_ = f.Init(ctx)
	_, _ = f.Update(ctx, runeKey('a'))

	if ctx.Config.Filter != nil {
		t.Fatal("expected no filter commit in manual mode without enter")
	}
}

func TestFilterToggleLiveSwitchesBehavior(t *testing.T) {
	ctx := newCtx(t)
	ctx.Config.LiveFilter = false

	f := New(sentinel)
	_ = f.Init(ctx)
	_, _ = f.Update(ctx, tea.KeyPressMsg(tea.Key{Code: 'l', Mod: tea.ModCtrl}))
	if !ctx.Config.LiveFilter {
		t.Fatal("expected live filter to be enabled after ctrl+l")
	}
}

func TestFilterViewReturnsNonEmpty(t *testing.T) {
	ctx := newCtx(t)
	ctx.Config.LiveFilter = true
	f := New(sentinel)
	_ = f.Init(ctx)
	v := f.Overlay(ctx)
	if v == "" {
		t.Fatal("expected non-empty view")
	}
}

func TestFilterHelpReturnsBindings(t *testing.T) {
	f := New(sentinel)
	bindings := f.Help()
	if bindings == nil {
		t.Fatal("expected non-nil help")
	}
	if len(bindings.ShortHelp()) == 0 {
		t.Fatal("expected non-empty help")
	}
}

func TestFilterConfirmReturnsToTarget(t *testing.T) {
	ctx := newCtx(t)
	f := New(sentinel)
	_ = f.Init(ctx)
	next, _ := f.Update(ctx, specialKey(tea.KeyEnter))
	if next != sentinel {
		t.Fatalf("expected returnTo on enter, got %T", next)
	}
}

func TestFilterCancelReturnsToTarget(t *testing.T) {
	ctx := newCtx(t)
	f := New(sentinel)
	_ = f.Init(ctx)
	next, _ := f.Update(ctx, specialKey(tea.KeyEscape))
	if next != sentinel {
		t.Fatalf("expected returnTo on esc, got %T", next)
	}
}
