package explorer

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/styles"
	stexmodel "github.com/SolracHQ/stex/internal/model"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

// buildTestTree creates a small tree in a temp dir, runs the
// async scan synchronously, and returns the scanned root.
func buildTestTree(t *testing.T) *stexmodel.Dir {
	t.Helper()
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "a.txt"), []byte("aaa"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "b.txt"), []byte("bbbbbb"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(tmp, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "sub", "c.txt"), []byte("cc"), 0o644); err != nil {
		t.Fatal(err)
	}
	state := &stexmodel.ScanState{}
	stexmodel.BuildTree(tmp, state)
	if state.Result == nil {
		t.Fatal("scan produced no root")
	}
	return state.Result
}

// newCtx builds a Context wired with the test tree and the
// standard table styles. The result is ready for Rebuild.
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
	core.Rebuild(ctx)
	return ctx
}

func runeKey(r rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: r, Text: string(r)})
}

func specialKey(code rune) tea.KeyPressMsg {
	return tea.KeyPressMsg(tea.Key{Code: code})
}

func TestExplorerRebuildPopulatesItems(t *testing.T) {
	ctx := newCtx(t)
	if len(ctx.Items) == 0 {
		t.Fatal("expected items after rebuild")
	}
	if ctx.Current != ctx.Root {
		t.Fatal("expected Current to be root")
	}
}

func TestExplorerStaysOnUnmatchedKey(t *testing.T) {
	ctx := newCtx(t)
	e := Explorer{}
	next, _ := e.Update(ctx, runeKey('z'))
	if next != nil {
		t.Fatalf("expected nil next (stay), got %T", next)
	}
}

func TestExplorerHelpToggleChangesShowAll(t *testing.T) {
	ctx := newCtx(t)
	e := Explorer{}
	prev := ctx.Help.ShowAll
	_, _ = e.Update(ctx, runeKey('?'))
	if ctx.Help.ShowAll == prev {
		t.Fatal("expected help toggle to flip ShowAll")
	}
}

func TestExplorerGroupOpensPicker(t *testing.T) {
	ctx := newCtx(t)
	e := Explorer{}
	next, _ := e.Update(ctx, runeKey('g'))
	if next == nil {
		t.Fatal("expected grouping picker on g, got nil")
	}
}

func TestExplorerClearFilterNoopWhenNoFilter(t *testing.T) {
	ctx := newCtx(t)
	e := Explorer{}
	if ctx.Config.Filter != nil {
		t.Fatal("test setup error: expected nil filter")
	}
	before := len(ctx.Items)
	_, _ = e.Update(ctx, runeKey('c'))
	if len(ctx.Items) != before {
		t.Fatal("expected no item change when no filter active")
	}
}

func TestExplorerClearFilterResetsItems(t *testing.T) {
	ctx := newCtx(t)
	e := Explorer{}
	re := regexp.MustCompile("a")
	ctx.Config.Filter = re
	core.Rebuild(ctx)
	_, _ = e.Update(ctx, runeKey('c'))
	if ctx.Config.Filter != nil {
		t.Fatal("expected filter to be cleared")
	}
}

func TestExplorerSearchOpensFilter(t *testing.T) {
	ctx := newCtx(t)
	e := Explorer{}
	next, _ := e.Update(ctx, runeKey('/'))
	if next == nil {
		t.Fatal("expected filter mode on /, got nil")
	}
}

func TestExplorerEnterOpensDirectory(t *testing.T) {
	ctx := newCtx(t)
	e := Explorer{}

	for i, item := range ctx.Items {
		if _, ok := item.(*stexmodel.Dir); ok {
			ctx.Table.SetCursor(i)
			break
		}
	}
	prev := ctx.Current

	_, _ = e.Update(ctx, specialKey(tea.KeyEnter))

	if ctx.Current == prev {
		t.Fatal("expected to navigate into directory on enter")
	}
}

func TestExplorerGoToParent(t *testing.T) {
	ctx := newCtx(t)
	e := Explorer{}

	for i, item := range ctx.Items {
		if _, ok := item.(*stexmodel.Dir); ok {
			ctx.Table.SetCursor(i)
			break
		}
	}
	_, _ = e.Update(ctx, specialKey(tea.KeyEnter))
	if ctx.Current == ctx.Root {
		t.Fatal("expected to enter subdir first")
	}
	_, _ = e.Update(ctx, specialKey(tea.KeyEscape))
	if ctx.Current != ctx.Root {
		t.Fatal("expected to return to root on esc")
	}
}

func TestExplorerUpdateInfoPopulatesOnFirstCall(t *testing.T) {
	ctx := newCtx(t)
	ctx.Info.Path = ""
	core.UpdateInfo(ctx)
	if ctx.Info.Path == "" {
		t.Fatal("expected Info.Path to be set after UpdateInfo")
	}
}

func TestExplorerUpdateInfoSkipsWhenPathUnchanged(t *testing.T) {
	ctx := newCtx(t)
	core.UpdateInfo(ctx)
	first := ctx.Info.Content
	core.UpdateInfo(ctx)
	if ctx.Info.Content != first {
		t.Fatal("expected Info.Content to be cached when cursor hasn't moved")
	}
}

func TestExplorerOverlayReturnsEmpty(t *testing.T) {
	ctx := newCtx(t)
	e := Explorer{}
	if v := e.Overlay(ctx); v != "" {
		t.Fatalf("expected empty view, got %q", v)
	}
}

func TestExplorerHelpReturnsBindings(t *testing.T) {
	e := Explorer{}
	bindings := e.Help()
	if bindings == nil {
		t.Fatal("expected non-nil help")
	}
	if len(bindings.ShortHelp()) == 0 {
		t.Fatal("expected non-empty short help")
	}
	if len(bindings.FullHelp()) == 0 {
		t.Fatal("expected non-empty full help")
	}
}

func TestExplorerFullHelpCoversAllModes(t *testing.T) {
	e := Explorer{}
	bindings := e.Help()

	wantKeys := []string{
		"c",      // clear filter
		"/",      // filter
		"ctrl+c", // quit
	}

	collected := collectKeyStrings(bindings.FullHelp())
	for _, k := range wantKeys {
		if !strings.Contains(collected, k) {
			t.Errorf("full help missing key %q\n%s", k, collected)
		}
	}
}

func collectKeyStrings(rows [][]key.Binding) string {
	var out []string
	for _, row := range rows {
		for _, b := range row {
			out = append(out, b.Help().Key)
		}
	}
	return strings.Join(out, " ")
}

func TestExplorerResizeTriggersRebuild(t *testing.T) {
	ctx := newCtx(t)
	e := Explorer{}
	prevTableWidth := ctx.Table.Width()
	// The app updates ctx.Width/Height before dispatching the
	// WindowSizeMsg to the mode. Simulate that here.
	ctx.Width = 200
	ctx.Height = 50
	_, _ = e.Update(ctx, tea.WindowSizeMsg{Width: 200, Height: 50})
	if ctx.Table.Width() == prevTableWidth {
		t.Fatal("expected table width to change on resize")
	}
}
