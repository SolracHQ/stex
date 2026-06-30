package scanning

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/explorer"
	"github.com/SolracHQ/stex/internal/model"

	tea "charm.land/bubbletea/v2"
)

func TestLoadingStartsInProgress(t *testing.T) {
	l := New(".")
	defer func() { l.state.Cancelled = true }()

	if l.state.Done {
		t.Fatal("expected Done=false on fresh Loading")
	}
	if l.state.Result != nil {
		t.Fatal("expected Result=nil on fresh Loading")
	}
}

func TestLoadingStaysOnTickWhenNotDone(t *testing.T) {
	l := &Loading{state: &model.ScanState{}}
	ctx := &core.Context{Width: 80, Height: 24}

	next, cmd := l.Update(ctx, scanTickMsg{})
	if next != nil {
		t.Fatalf("expected nil next, got %T", next)
	}
	if cmd == nil {
		t.Fatal("expected tick command to be returned")
	}
}

func TestLoadingTransitionsWhenDone(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "f.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	state := &model.ScanState{}
	model.BuildTree(tmp, state)
	if state.Result == nil {
		t.Fatal("scan produced no root")
	}
	l := &Loading{state: state}
	ctx := &core.Context{Width: 80, Height: 24}

	next, cmd := l.Update(ctx, scanTickMsg{})
	if cmd != nil {
		t.Fatalf("expected nil cmd on transition, got %v", cmd)
	}
	if _, ok := next.(*explorer.Explorer); !ok {
		t.Fatalf("expected *Explorer, got %T", next)
	}
	if ctx.Root != state.Result {
		t.Fatal("expected ctx.Root to be set")
	}
	if ctx.Current != state.Result {
		t.Fatal("expected ctx.Current to be set")
	}
	if !ctx.Ready {
		t.Fatal("expected ctx.Ready to be true")
	}
}

func TestLoadingViewEmptyWhenNoDims(t *testing.T) {
	l := &Loading{state: &model.ScanState{}}
	ctx := &core.Context{Width: 0, Height: 0}
	if v := l.View(ctx); v != "" {
		t.Fatalf("expected empty view, got %q", v)
	}
}

func TestLoadingViewRendersOnDims(t *testing.T) {
	l := &Loading{state: &model.ScanState{}}
	ctx := &core.Context{Width: 80, Height: 24}
	v := l.View(ctx)
	if v == "" {
		t.Fatal("expected non-empty view")
	}
}

func TestLoadingHelpIsNil(t *testing.T) {
	l := &Loading{state: &model.ScanState{}}
	if l.Help() != nil {
		t.Fatal("expected nil help")
	}
}

func TestTickReturnsMsg(t *testing.T) {
	cmd := tick()
	if cmd == nil {
		t.Fatal("expected non-nil tick command")
	}
	done := make(chan tea.Msg, 1)
	go func() { done <- cmd() }()
	select {
	case msg := <-done:
		if _, ok := msg.(scanTickMsg); !ok {
			t.Fatalf("expected scanTickMsg, got %T", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("tick did not return within 1s")
	}
}
