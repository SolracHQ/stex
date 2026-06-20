package tui

import (
	"fmt"
	"math"
	"strings"
	"time"

	stexmodel "github.com/SolracHQ/stex/model"
	stexview "github.com/SolracHQ/stex/view"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
)

type screen int

const (
	screenMain screen = iota
	screenHelp
)

type scanTickMsg struct{}

type model struct {
	path string

	root    *stexmodel.Dir
	current *stexmodel.Dir
	cfg     stexmodel.Config
	items   []stexmodel.TreeItem

	ready  bool
	screen screen

	scanState *stexmodel.ScanState

	width  int
	height int
	keys   keyMap
	help   help.Model
	tableView    table.Model
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

	tableView := table.New(
		table.WithColumns([]table.Column{
			{Title: " Size%", Width: 9},
			{Title: "Size", Width: 12},
			{Title: "Name", Width: 40},
		}),
		table.WithFocused(true),
		table.WithStyles(tableStyles()),
	)

	return &model{
		path:   path,
		cfg:    stexmodel.DefaultConfig(),
		width:  80,
		height: 24,
		keys:   appKeys,
		help:   helpModel,
		tableView:    tableView,
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

// rebuildRows refreshes the table columns (with sort indicators) and rows
// from the current items slice.
func (m *model) rebuildRows() {
	cols := m.tableView.Columns()
	nameWidth := 40
	if len(cols) > 2 {
		nameWidth = cols[2].Width
	}

	var sizeLabel, nameLabel string
	switch m.cfg.SortBy {
	case stexmodel.SortBySize:
		if m.cfg.SortOrder == stexmodel.Descending {
			sizeLabel = "Size↓"
		} else {
			sizeLabel = "Size↑"
		}
		nameLabel = "Name"
	case stexmodel.SortByName:
		sizeLabel = "Size"
		if m.cfg.SortOrder == stexmodel.Descending {
			nameLabel = "Name↓"
		} else {
			nameLabel = "Name↑"
		}
	}

	m.tableView.SetColumns([]table.Column{
		{Title: " Size%", Width: 9},
		{Title: sizeLabel, Width: 12},
		{Title: nameLabel, Width: nameWidth},
	})

	rows := make([]table.Row, 0, len(m.items))
	for _, item := range m.items {
		rows = append(rows, itemToRow(item, m.cfg.IconStyle))
	}
	m.tableView.SetRows(rows)
}

// itemToRow converts a TreeItem into a table row with a gradient-coloured
// percentage cell (green→yellow→red based on disk usage). The gradient ANSI
// closes with \033[39m so it only resets foreground, not the selection
// background. In emoji mode, files and directories get a prefix emoji.
func itemToRow(item stexmodel.TreeItem, iconStyle stexmodel.IconStyle) table.Row {
	switch item.Kind {
	case stexmodel.TKFile:
		percent := item.File.Size.PercentOf(item.File.Parent.Size())
		gradientCode := gradientANSI(percent)
		name := item.File.Name
		if iconStyle == stexmodel.IconEmoji {
			name = "📄 " + name
		}
		return table.Row{
			gradientCode + fmt.Sprintf("%5.2f%%", percent) + "\033[39m",
			gradientCode + " " + item.File.Size.String() + " \033[39m",
			" " + name + " ",
		}
	case stexmodel.TKDir:
		percent := item.Dir.Size().PercentOf(item.Dir.Parent.Size())
		gradientCode := gradientANSI(percent)
		name := item.Dir.Name
		if iconStyle == stexmodel.IconEmoji {
			name = "📁 " + name
		}
		return table.Row{
			gradientCode + fmt.Sprintf("%5.2f%%", percent) + "\033[39m",
			gradientCode + " " + item.Dir.Size().String() + " \033[39m",
			" " + name + " ",
		}
	case stexmodel.TKUpLink:
		return table.Row{"", "", "   ..  "}
	}
	return table.Row{}
}

// gradientANSI returns an ANSI foreground code for a percentage 0–100,
// mapping green → yellow → red.
func gradientANSI(percent float64) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	var r, g int
	if percent <= 50 {
		factor := percent / 50.0
		r = int(math.Round(255 * factor))
		g = 255
	} else {
		factor := (percent - 50) / 50.0
		r = 255
		g = int(math.Round(255 * (1 - factor)))
	}
	return fmt.Sprintf("\033[38;2;%d;%d;0m", r, g)
}

// Update handles all messages: window resize, scan progress ticks, key
// presses, mouse clicks and scrolls.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg), nil

	case scanTickMsg:
		return m.handleScanTick()

	case tea.KeyPressMsg:
		return m.handleKeyPress(msg)

	case tea.MouseClickMsg:
		return m.handleMouseClick(msg)

	case tea.MouseWheelMsg:
		return m.handleMouseWheel(msg), nil
	}

	return m, nil
}

// handleWindowSize updates the terminal dimensions and resizes the table.
func (m *model) handleWindowSize(msg tea.WindowSizeMsg) *model {
	m.width = msg.Width
	m.height = msg.Height
	m.help.SetWidth(msg.Width - 4)

	innerWidth := msg.Width - 2
	innerHeight := msg.Height - 2
	nameWidth := innerWidth - 9 - 12 - 2
	if nameWidth < 10 {
		nameWidth = 10
	}
	cols := m.tableView.Columns()
	if len(cols) > 2 {
		cols[2].Width = nameWidth
		m.tableView.SetColumns(cols)
	}
	m.tableView.SetWidth(innerWidth)
	m.tableView.SetHeight(innerHeight - 3)
	return m
}

// handleScanTick polls the scan state. When the scan finishes it wires the
// root directory and switches to the main view.
func (m *model) handleScanTick() (tea.Model, tea.Cmd) {
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
		m.rebuildRows()
		return m, nil
	}
	return m, m.scanTick()
}

// handleKeyPress routes keyboard input. It delegates navigation to the table,
// handles quit/help/enter/back, and toggles sort, group, order, and icons.
func (m *model) handleKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	if !m.ready {
		if key.Matches(msg, m.keys.Quit) {
			m.scanState.Mu.Lock()
			m.scanState.Cancelled = true
			m.scanState.Mu.Unlock()
			return m, tea.Quit
		}
		return m, nil
	}

	if m.screen == screenHelp {
		if key.Matches(msg, m.keys.Help) || key.Matches(msg, m.keys.Back) {
			m.screen = screenMain
		}
		return m, nil
	}

	if key.Matches(msg, m.keys.Quit) {
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.tableView, cmd = m.tableView.Update(msg)
	if cmd != nil {
		return m, cmd
	}

	switch {
	case key.Matches(msg, m.keys.Help):
		m.screen = screenHelp
		m.help.ShowAll = true

	case key.Matches(msg, m.keys.Enter):
		m.enterSelected()

	case key.Matches(msg, m.keys.Back):
		m.goToParent()

	case key.Matches(msg, m.keys.Sort):
		m.cfg.SortBy.Toggle()
		m.rebuildItems()

	case key.Matches(msg, m.keys.Group):
		m.cfg.Grouping.Toggle()
		m.rebuildItems()

	case key.Matches(msg, m.keys.Order):
		m.cfg.SortOrder.Toggle()
		m.rebuildItems()

	case key.Matches(msg, m.keys.Icons):
		m.cfg.IconStyle.Toggle()
		m.rebuildRows()
	}

	return m, nil
}

// enterSelected opens the currently selected directory or navigates up.
func (m *model) enterSelected() {
	if len(m.items) == 0 {
		return
	}
	idx := m.tableView.Cursor()
	if idx < 0 || idx >= len(m.items) {
		return
	}
	item := m.items[idx]
	switch item.Kind {
	case stexmodel.TKDir:
		m.current = item.Dir
		m.items = stexmodel.ComputeItems(m.current, m.cfg)
		m.rebuildRows()
		m.tableView.SetCursor(0)
	case stexmodel.TKUpLink:
		m.current = item.Parent
		m.items = stexmodel.ComputeItems(m.current, m.cfg)
		m.rebuildRows()
		m.tableView.SetCursor(0)
	}
}

// goToParent moves the view up one directory level.
func (m *model) goToParent() {
	if m.current.Parent != nil {
		m.current = m.current.Parent
		m.rebuildItems()
		m.tableView.SetCursor(0)
	}
}

// rebuildItems recomputes the item list from the current directory and
// refreshes the table rows.
func (m *model) rebuildItems() {
	m.items = stexmodel.ComputeItems(m.current, m.cfg)
	m.rebuildRows()
}

// handleMouseClick selects the clicked data row and opens directories.
// Clicks on the table header row are forwarded to handleHeaderClick to
// trigger sorting.
func (m *model) handleMouseClick(msg tea.MouseClickMsg) (tea.Model, tea.Cmd) {
	if !m.ready || m.screen != screenMain {
		return m, nil
	}
	mouse := msg.Mouse()
	clickY := mouse.Y - 1
	clickX := mouse.X - 1

	if clickY == 2 {
		m.handleHeaderClick(clickX)
		return m, nil
	}

	rowIndex := clickY - 3
	if rowIndex < 0 || rowIndex >= len(m.items) {
		return m, nil
	}
	m.tableView.SetCursor(rowIndex)

	item := m.items[rowIndex]
	switch item.Kind {
	case stexmodel.TKDir:
		m.current = item.Dir
		m.rebuildItems()
		m.tableView.SetCursor(0)
	case stexmodel.TKUpLink:
		m.current = item.Parent
		m.rebuildItems()
		m.tableView.SetCursor(0)
	}
	return m, nil
}

// handleHeaderClick sorts by the clicked column or toggles the sort order
// when the column is already active.
func (m *model) handleHeaderClick(clickX int) {
	cols := m.tableView.Columns()
	x := 1
	for i, col := range cols {
		columnWidth := col.Width
		if clickX >= x && clickX < x+columnWidth {
			switch i {
			case 1: // Size column
				if m.cfg.SortBy == stexmodel.SortBySize {
					m.cfg.SortOrder.Toggle()
				} else {
					m.cfg.SortBy = stexmodel.SortBySize
				}
				m.rebuildItems()
			case 2: // Name column
				if m.cfg.SortBy == stexmodel.SortByName {
					m.cfg.SortOrder.Toggle()
				} else {
					m.cfg.SortBy = stexmodel.SortByName
				}
				m.rebuildItems()
			}
			break
		}
		x += columnWidth
	}
}

// handleMouseWheel scrolls the table up or down by three rows.
func (m *model) handleMouseWheel(msg tea.MouseWheelMsg) *model {
	if !m.ready || m.screen != screenMain {
		return m
	}
	if msg.Button == tea.MouseWheelUp {
		m.tableView.MoveUp(3)
	} else {
		m.tableView.MoveDown(3)
	}
	return m
}

// View renders the current screen: scan progress, main file listing with
// table, or the main listing overlaid with the help dialog.
func (m *model) View() tea.View {
	if !m.ready {
		return m.loadingView()
	}

	innerWidth := m.width - 2
	innerHeight := m.height - 2

	lines := make([]string, innerHeight)

	lines[0] = stexview.Title(m.current, innerWidth, m.cfg.IconStyle)
	lines[1] = strings.Repeat("─", innerWidth)

	tableContent := m.tableView.View()
	tableLines := strings.Split(tableContent, "\n")
	availableSlots := innerHeight - 2

	for i := 0; i < availableSlots-1 && i < len(tableLines); i++ {
		lines[2+i] = tableLines[i]
	}
	for i := len(tableLines); i < availableSlots-1; i++ {
		lines[2+i] = ""
	}

	grouping := stexmodel.GroupingString(m.cfg.Grouping)
	lines[innerHeight-1] = fmt.Sprintf(" (q) quit  (?) help  |  %s", grouping)

	inner := strings.Join(lines, "\n")

	if m.screen == screenHelp {
		inner = overlayCenter(inner, m.buildHelpDialog(innerWidth, innerHeight))
	}

	v := tea.NewView(borderStyle.Render(inner))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

// loadingView renders the scan progress dialog centred in the terminal.
func (m *model) loadingView() tea.View {
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
	v := tea.NewView(centered)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeAllMotion
	return v
}

// buildHelpDialog returns a bordered box containing key bindings and the
// current configuration summary.
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
		dialogLines = append(dialogLines, dialogFooterStyle.Render(" (?) close  |  (q) quit "))

	content := strings.Join(dialogLines, "\n")
	return dialogBoxStyle.Render(content)
}
