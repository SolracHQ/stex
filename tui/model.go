package tui

import (
	"fmt"
	"strings"
	"time"

	stexmodel "github.com/SolracHQ/stex/model"
	stexview "github.com/SolracHQ/stex/view"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

// screen identifies the currently active view.
type screen int

const (
	screenMain screen = iota // file listing
	screenHelp               // help overlay
)

// scanCompleteMsg is sent when the async scan finishes. Kept for future use.
type scanCompleteMsg struct {
	root *stexmodel.Dir
	err  error
}

// scanTickMsg is sent periodically while the scan is in progress.
type scanTickMsg struct{}

// model is the top-level Bubble Tea model. It holds all application state.
type model struct {
	path string

	root    *stexmodel.Dir
	current *stexmodel.Dir
	cfg     stexmodel.Config
	items   []stexmodel.TreeItem

	ready  bool
	cursor int
	err    error
	screen screen

	scanState *stexmodel.ScanState

	width  int
	height int
	keys   keyMap
	help   help.Model
}

// New creates and returns a new Bubble Tea model for the given directory path.
func New(path string) tea.Model {
	helpModel := help.New()
	helpModel.Styles.FullKey = helpKeyStyle
	helpModel.Styles.FullDesc = helpDescStyle
	helpModel.Styles.FullSeparator = helpSepStyle
	helpModel.Styles.ShortKey = helpKeyStyle
	helpModel.Styles.ShortDesc = helpDescStyle
	helpModel.Styles.ShortSeparator = helpSepStyle

	return &model{
		path:   path,
		cfg:    stexmodel.DefaultConfig(),
		width:  80,
		height: 24,
		keys:   appKeys,
		help:   helpModel,
	}
}

// Init starts the asynchronous directory scan.
func (m *model) Init() tea.Cmd {
	m.scanState = &stexmodel.ScanState{}

	go func() {
		stexmodel.BuildTree(m.path, m.scanState)
	}()

	return m.scanTick()
}

// scanTick returns a tea.Cmd that fires scanTickMsg after 80ms.
func (m *model) scanTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return scanTickMsg{}
	})
}

// Update handles all messages, including window resize, scan progress ticks,
// and key presses across both main and help screens.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.SetWidth(msg.Width - 4)
		return m, nil

	case scanTickMsg:
		if m.ready {
			return m, nil
		}

		m.scanState.Mu.Lock()
		done := m.scanState.Done
		root := m.scanState.Result
		m.scanState.Mu.Unlock()

		if done {
			m.root = root
			m.current = root
			m.ready = true
			m.items = stexmodel.ComputeItems(m.current, m.cfg)
			return m, nil
		}

		return m, m.scanTick()

	case scanCompleteMsg:
		return m, nil

	case tea.KeyPressMsg:
		if key.Matches(msg, m.keys.Quit) {
			if !m.ready {
				m.scanState.Mu.Lock()
				m.scanState.Cancelled = true
				m.scanState.Mu.Unlock()
			}
			return m, tea.Quit
		}
		if !m.ready {
			return m, nil
		}

		if m.screen == screenHelp {
			if key.Matches(msg, m.keys.Help) || key.Matches(msg, m.keys.Back) {
				m.screen = screenMain
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, m.keys.Help):
			m.screen = screenHelp
			m.help.ShowAll = true

		case key.Matches(msg, m.keys.Up):
			if len(m.items) == 0 {
				return m, nil
			}
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.items) - 1
			}

		case key.Matches(msg, m.keys.Down):
			if len(m.items) == 0 {
				return m, nil
			}
			m.cursor++
			if m.cursor >= len(m.items) {
				m.cursor = 0
			}

		case key.Matches(msg, m.keys.Enter):
			if len(m.items) == 0 {
				return m, nil
			}
			item := m.items[m.cursor]
			switch item.Kind {
			case stexmodel.TKDir:
				m.current = item.Dir
				m.cursor = 0
				m.items = stexmodel.ComputeItems(m.current, m.cfg)
			case stexmodel.TKUpLink:
				m.current = item.Parent
				m.cursor = 0
				m.items = stexmodel.ComputeItems(m.current, m.cfg)
			}

		case key.Matches(msg, m.keys.Back):
			if m.current.Parent != nil {
				m.current = m.current.Parent
				m.cursor = 0
				m.items = stexmodel.ComputeItems(m.current, m.cfg)
			}

		case key.Matches(msg, m.keys.Sort):
			m.cfg.SortBy.Toggle()
			m.items = stexmodel.ComputeItems(m.current, m.cfg)

		case key.Matches(msg, m.keys.Group):
			m.cfg.Grouping.Toggle()
			m.items = stexmodel.ComputeItems(m.current, m.cfg)

		case key.Matches(msg, m.keys.Order):
			m.cfg.SortOrder.Toggle()
			m.items = stexmodel.ComputeItems(m.current, m.cfg)
		}
	}

	return m, nil
}

// View renders the current screen. It shows the scan progress dialog when not
// ready, the main file listing otherwise, optionally overlaid with the help
// dialog.
func (m *model) View() tea.View {
	if !m.ready {
		innerWidth := m.width - 2
		progress := stexview.Progress(m.scanState, innerWidth)
		dialog := scanBorderStyle.Render(progress)

		emptyLine := strings.Repeat(" ", m.width)
		bgLines := make([]string, m.height)
		for line := range bgLines {
			bgLines[line] = emptyLine
		}
		bg := strings.Join(bgLines, "\n")
		centered := overlayCenter(bg, dialog)

		view := tea.NewView(centered)
		view.AltScreen = true
		return view
	}

	innerWidth := m.width - 2
	innerHeight := m.height - 2

	lines := make([]string, innerHeight)

	lines[0] = stexview.Title(m.current, innerWidth)
	lines[1] = strings.Repeat("─", innerWidth)
	lines[2] = stexview.ColumnHeaders(m.cfg, innerWidth)

	itemHeight := innerHeight - 4
	items := stexview.List(m.items, m.cursor, itemHeight, innerWidth, stexview.DefaultDisplay)
	itemLines := strings.Split(items, "\n")
	for line := 0; line < itemHeight; line++ {
		if line < len(itemLines) {
			lines[3+line] = itemLines[line]
		} else {
			lines[3+line] = ""
		}
	}

	grouping := stexmodel.GroupingString(m.cfg.Grouping)
	lines[innerHeight-1] = fmt.Sprintf(" (q) quit  (h/?) help  |  %s", grouping)

	inner := strings.Join(lines, "\n")

	if m.screen == screenHelp {
		inner = overlayCenter(inner, m.buildHelpDialog(innerWidth, innerHeight))
	}

	view := tea.NewView(borderStyle.Render(inner))
	view.AltScreen = true
	return view
}

// buildHelpDialog returns a bordered dialog box containing the key binding
// listing and current configuration summary.
func (m *model) buildHelpDialog(innerWidth, innerHeight int) string {
	helpContent := m.help.View(m.keys)

	configLines := fmt.Sprintf("sort by: %s %s  |  grouping: %s",
		map[bool]string{true: "name", false: "size"}[m.cfg.SortBy == stexmodel.SortByName],
		map[bool]string{true: "↑", false: "↓"}[m.cfg.SortOrder == stexmodel.Ascending],
		stexmodel.GroupingString(m.cfg.Grouping),
	)

	helpLines := strings.Split(helpContent, "\n")

	dialogLines := make([]string, 0, len(helpLines)+6)
	dialogLines = append(dialogLines, dialogTitleStyle.Render(" Help "))
	dialogLines = append(dialogLines, "")
	dialogLines = append(dialogLines, helpLines...)
	dialogLines = append(dialogLines, "")
	dialogLines = append(dialogLines, dialogFooterStyle.Render(configLines))
	dialogLines = append(dialogLines, "")
	dialogLines = append(dialogLines, dialogFooterStyle.Render(" (?/h) close  |  (q) quit "))

	content := strings.Join(dialogLines, "\n")
	return dialogBoxStyle.Render(content)
}
