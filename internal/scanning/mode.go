// Package scanning implements the loading mode shown while the initial directory tree is
// being scanned. It owns the ScanState (built by model.BuildTree in a goroutine) and polls it
// on a tick until completion, then transitions to the explorer.
package scanning

import (
	"time"

	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/explorer"
	"github.com/SolracHQ/stex/internal/model"

	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
)

// scanTickInterval is the period between progress polls. Short enough that the dialog feels
// live, long enough that the UI thread is not starved.
const scanTickInterval = 80 * time.Millisecond

// scanTickMsg is the tick that drives the progress polling. The concrete type is private so
// other packages cannot emit it.
type scanTickMsg struct{}

// Loading is the mode shown while the initial directory tree is being scanned. It owns its
// own ScanState, not the shared Context, because no other mode needs the scan progress data.
type Loading struct {
	state *model.ScanState
}

// New starts a background scan of path and returns the Loading mode bound to its ScanState.
// The caller is expected to install the mode as the initial active mode.
func New(path string) *Loading {
	state := &model.ScanState{}
	go model.BuildTree(path, state)
	return &Loading{state: state}
}

// Init returns the first poll command. The Bubble Tea runtime will run it and the resulting
// message will be delivered to Update, which schedules the next poll.
func (l *Loading) Init(_ *core.Context) tea.Cmd {
	return tick()
}

// Update handles the scan tick. When the scan is done it sets the Context's root, current
// directory, and Ready flag, and returns a new explorer mode. While the scan is still
// running it schedules the next poll.
func (l *Loading) Update(ctx *core.Context, msg tea.Msg) (core.Mode, tea.Cmd) {
	switch msg.(type) {
	case scanTickMsg:
		l.state.Mu.Lock()
		done := l.state.Done
		root := l.state.Result
		l.state.Mu.Unlock()
		if done {
			ctx.Root = root
			ctx.Current = root
			ctx.Ready = true
			return &explorer.Explorer{}, nil
		}
		return nil, tick()
	}
	return nil, nil
}

// View returns the progress dialog rendered at the context's dimensions. Returns "" when
// dimensions are not yet known, so the program does not flash an empty box on the first frame.
func (l *Loading) View(ctx *core.Context) string {
	if ctx.Width == 0 || ctx.Height == 0 {
		return ""
	}
	return progressBox(l.state, ctx.Width, ctx.Height)
}

// Help returns nil, the loading screen has no key bindings.
func (l *Loading) Help() help.KeyMap { return nil }

// tick returns a tea.Cmd that emits a scanTickMsg after one scanTickInterval. The Bubble Tea
// runtime will deliver the message to the active mode's Update, which decides whether to
// schedule another tick.
func tick() tea.Cmd {
	return tea.Tick(scanTickInterval, func(_ time.Time) tea.Msg {
		return scanTickMsg{}
	})
}
