package scanning

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/model"

	tea "charm.land/bubbletea/v2"
)

func TestLoadingStartsInProgress(t *testing.T) {
	l := New(".", config.DefaultConfig())
	defer func() { l.state.Cancelled = true }()

	if l.state.Done {
		t.Fatal("expected Done=false on fresh Loading")
	}
	if l.state.Result != nil {
		t.Fatal("expected Result=nil on fresh Loading")
	}
}

func TestLoadingInitReturnsTick(t *testing.T) {
	l := &Loading{state: &model.ScanState{}}
	cmd := l.Init()
	if cmd == nil {
		t.Fatal("expected tick command from Init")
	}
}

func TestLoadingReschedulesTickWhileScanning(t *testing.T) {
	l := &Loading{state: &model.ScanState{}}
	next, cmd := l.Update(scanTickMsg{})
	if cmd == nil {
		t.Fatal("expected tick command")
	}
	if next != l {
		t.Fatal("expected self as next model when scan is running")
	}
}

func TestLoadingReturnsAppOnScanComplete(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "f.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	state := &model.ScanState{}
	model.BuildTree(tmp, state)
	if state.Result == nil {
		t.Fatal("scan produced no root")
	}
	l := &Loading{path: tmp, cfg: config.DefaultConfig(), state: state}

	next, cmd := l.Update(scanTickMsg{})
	if cmd != nil {
		t.Fatalf("expected nil cmd on transition, got %v", cmd)
	}
	if next == nil {
		t.Fatal("expected a tea.Model on scan complete")
	}
}

func TestLoadingViewBlankWhenNoDims(t *testing.T) {
	l := &Loading{state: &model.ScanState{}}
	v := l.View()
	if v.Content != "" {
		t.Fatalf("expected empty view, got %q", v.Content)
	}
}

func TestLoadingViewRendersAfterResize(t *testing.T) {
	l := &Loading{state: &model.ScanState{}}
	l.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	v := l.View()
	if v.Content == "" {
		t.Fatal("expected non-empty view after resize")
	}
}

func TestLoadingStoresWindowSize(t *testing.T) {
	l := &Loading{state: &model.ScanState{}}
	l.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	if l.width != 120 || l.height != 40 {
		t.Fatalf("expected 120x40, got %dx%d", l.width, l.height)
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
