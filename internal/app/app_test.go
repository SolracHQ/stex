package app

import (
	"charm.land/bubbles/v2/help"
	"testing"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/scanning"

	tea "charm.land/bubbletea/v2"
)

func TestAppStartsInScanning(t *testing.T) {
	m := New(".", config.DefaultConfig())
	a, ok := m.(*App)
	if !ok {
		t.Fatalf("expected *App, got %T", m)
	}
	if a.ctx.Ready {
		t.Fatal("expected Ready=false at start")
	}
	if a.mode == nil {
		t.Fatal("expected a mode to be set")
	}
	if _, ok := a.mode.(*scanning.Loading); !ok {
		t.Fatalf("expected *scanning.Loading, got %T", a.mode)
	}
}

func TestAppTableFocusedAtStart(t *testing.T) {
	m := New(".", config.DefaultConfig())
	a := m.(*App)
	if !a.ctx.Table.Focused() {
		t.Fatal("expected table to be focused at start so key navigation works")
	}
}

func TestAppHandlesResizeBeforeInit(t *testing.T) {
	a := &App{
		ctx:  &core.Context{Width: 80, Height: 24, Keys: core.DefaultKeys()},
		mode: fakeNext{},
	}
	_, _ = a.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if a.ctx.Width != 120 || a.ctx.Height != 40 {
		t.Fatalf("expected dims 120x40, got %dx%d", a.ctx.Width, a.ctx.Height)
	}
}

// fakeMode is a test mode that returns next from Update, used to verify the app installs the
// returned mode and runs its Init.
type fakeMode struct {
	next core.Mode
}

func (f *fakeMode) Init(_ *core.Context) tea.Cmd { return nil }

func (f *fakeMode) Update(_ *core.Context, _ tea.Msg) (core.Mode, tea.Cmd) {
	return f.next, nil
}

func (f *fakeMode) View(_ *core.Context) string { return "" }

func (f *fakeMode) Help() help.KeyMap { return nil }

type fakeNext struct{}

func (fakeNext) Init(_ *core.Context) tea.Cmd { return nil }

func (fakeNext) Update(_ *core.Context, _ tea.Msg) (core.Mode, tea.Cmd) {
	return nil, nil
}

func (fakeNext) View(_ *core.Context) string { return "" }

func (fakeNext) Help() help.KeyMap { return nil }

func TestAppSwapsModeWhenReturned(t *testing.T) {
	a := &App{
		ctx:  &core.Context{Width: 80, Height: 24, Keys: core.DefaultKeys()},
		mode: &fakeMode{next: fakeNext{}},
	}
	_, _ = a.Update(nil)
	if _, ok := a.mode.(fakeNext); !ok {
		t.Fatalf("expected mode swap to fakeNext, got %T", a.mode)
	}
}

func TestOverlayCenterEmptyForeground(t *testing.T) {
	bg := "abc\ndef"
	if got := overlayCenter(bg, ""); got != bg {
		t.Fatalf("expected background returned, got %q", got)
	}
}

func TestOverlayCenterEmptyBackground(t *testing.T) {
	if got := overlayCenter("", "abc"); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}
