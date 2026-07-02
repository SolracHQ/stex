// Package scanning implements the initial loading model shown while the directory tree is
// being scanned. It is a standalone tea.Model, not a core.Mode, because no shared context
// exists yet. When the scan completes, Update returns the App model built from the scanned
// tree and Bubble Tea swaps the active model.
package scanning

import (
	"strings"
	"time"

	"github.com/SolracHQ/stex/internal/app"
	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/model"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// scanTickInterval is the period between progress polls. Short enough that the dialog feels
// live, long enough that the UI thread is not starved.
const scanTickInterval = 80 * time.Millisecond

// scanTickMsg is the tick that drives the progress polling. The concrete type is private so
// other packages cannot emit it.
type scanTickMsg struct{}

// Loading is the model shown while the initial directory tree is being scanned. It owns its
// own ScanState, built by BuildTree in a goroutine, and polls it on a tick until completion,
// then swaps itself for the App model. Window resize and the global quit key are handled here
// because no wrapping model exists.
type Loading struct {
	path   string
	cfg    config.Config
	state  *model.ScanState
	width  int
	height int
	quit   key.Binding
}

// New starts a background scan of path with the given resolved config and returns the Loading
// model bound to its ScanState. The caller passes it to tea.NewProgram.
func New(path string, cfg config.Config) *Loading {
	state := &model.ScanState{}
	go model.BuildTree(path, state)
	return &Loading{
		path:  path,
		cfg:   cfg,
		state: state,
		quit:  core.DefaultKeys().Quit,
	}
}

// Init returns the first poll command. The Bubble Tea runtime will run it and the resulting
// message will be delivered to Update, which schedules the next poll.
func (load *Loading) Init() tea.Cmd {
	return tick()
}

// Update polls the ScanState. On completion it returns the App model built from the scanned
// tree, Bubble Tea then swaps the active model. Window resize and the global quit key are
// intercepted here so the loading screen responds to resize and the user can abort the scan.
func (load *Loading) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		load.width = msg.Width
		load.height = msg.Height
		return load, nil
	case tea.KeyPressMsg:
		if key.Matches(msg, load.quit) {
			load.state.Mu.Lock()
			load.state.Cancelled = true
			load.state.Mu.Unlock()
			return load, tea.Quit
		}
	case scanTickMsg:
		load.state.Mu.Lock()
		done := load.state.Done
		root := load.state.Result
		load.state.Mu.Unlock()
		if done {
			cmd := emitSize(load.width, load.height)
			return app.New(load.path, load.cfg, root), cmd
		}
		return load, tick()
	}
	return load, nil
}

// View renders the progress dialog vertically centered, blank when dimensions are not yet
// known so the program does not flash an empty box on the first frame.
func (load *Loading) View() tea.View {
	if load.width == 0 || load.height == 0 {
		return tea.NewView("")
	}
	body := progressBox(load.state, load.width, load.height)
	lines := strings.Count(body, "\n") + 1
	pad := (load.height - lines) / 2
	if pad > 0 {
		body = strings.Repeat("\n", pad) + body
	}
	view := tea.NewView(body)
	view.AltScreen = true
	view.MouseMode = tea.MouseModeAllMotion
	return view
}

// tick returns a tea.Cmd that emits a scanTickMsg after one scanTickInterval. The Bubble Tea
// runtime will deliver the message to the active model's Update, which decides whether to
// schedule another tick.
func tick() tea.Cmd {
	return tea.Tick(scanTickInterval, func(_ time.Time) tea.Msg {
		return scanTickMsg{}
	})
}

// emitSize returns a cmd that sends a WindowSizeMsg with the given dimensions. When the
// loading model transitions to the App, this cmd ensures the App immediately knows the
// terminal size without waiting for the terminal to send its own resize message, which can
// result in a blank frame.
func emitSize(w, h int) tea.Cmd {
	if w == 0 || h == 0 {
		return nil
	}
	return func() tea.Msg {
		return tea.WindowSizeMsg{Width: w, Height: h}
	}
}
