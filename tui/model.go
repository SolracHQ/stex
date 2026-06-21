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
	"charm.land/lipgloss/v2"
	tea "charm.land/bubbletea/v2"
)

const (
	splitViewThreshold = 80
	scanTickInterval   = 80 * time.Millisecond
	scrollLines        = 3
	loadingMargin      = 6
	minDialogHeight    = 10
	maxChildren        = 10
	infoPanelWidth     = 80
	infoPanelHeight    = 10
)

type scanTickMsg struct{}

type model struct {
	path string

	root    *stexmodel.Dir
	current *stexmodel.Dir
	cfg     stexmodel.Config
	items   []stexmodel.TreeItem

	ready  bool

	scanState *stexmodel.ScanState

	width  int
	height int
	keys   keyMap
	help   help.Model
	tableView table.Model

	infoPath    string
	infoContent string
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

// scanTick returns a command that polls the async scan state.
func (m *model) scanTick() tea.Cmd {
	return tea.Tick(scanTickInterval, func(_ time.Time) tea.Msg {
		return scanTickMsg{}
	})
}

// rebuildRows refreshes the table columns and rows from the current items.
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
// percentage cell based on disk usage.
func itemToRow(item stexmodel.TreeItem, iconStyle stexmodel.IconStyle) table.Row {
	switch item.Kind {
	case stexmodel.TKFile, stexmodel.TKDir:
		return buildRow(item.Name(), item.Icon(), item.Size(), item.ParentDir().Size(), iconStyle)
	case stexmodel.TKUpLink:
		return table.Row{"", "", "   ..  "}
	}
	return table.Row{}
}

func buildRow(name, emoji string, size, parentSize stexmodel.Size, iconStyle stexmodel.IconStyle) table.Row {
	percent := size.PercentOf(parentSize)
	gradientCode := gradientANSI(percent)
	if iconStyle == stexmodel.IconEmoji {
		name = emoji + " " + name
	}
	return table.Row{
		gradientCode + fmt.Sprintf("%5.2f%%", percent) + "\033[39m",
		gradientCode + " " + size.String() + " \033[39m",
		" " + name + " ",
	}
}

// gradientANSI maps a percentage 0-100 to a green-yellow-red ANSI color.
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

// handleScanTick polls the scan state and switches to the main view when done.
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

	if key.Matches(msg, m.keys.Quit) {
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.tableView, cmd = m.tableView.Update(msg)
	if cmd != nil {
		return m, cmd
	}
	m.updateInfo()

	switch {
	case key.Matches(msg, m.keys.Help):
		m.help.ShowAll = !m.help.ShowAll

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
	m.infoPath = ""
	if len(m.items) == 0 {
		return
	}
	idx := m.tableView.Cursor()
	if idx < 0 || idx >= len(m.items) {
		return
	}
	m.current.LastCursor = idx

	item := m.items[idx]
	switch item.Kind {
	case stexmodel.TKDir:
		m.current = item.Dir
	case stexmodel.TKUpLink:
		m.current = item.Parent
	default:
		return
	}
	m.items = stexmodel.ComputeItems(m.current, m.cfg)
	m.rebuildRows()
	m.tableView.SetCursor(m.current.LastCursor)
}

// goToParent moves the view up one directory level.
func (m *model) goToParent() {
	m.infoPath = ""
	if m.current.Parent == nil {
		return
	}
	m.current.LastCursor = m.tableView.Cursor()
	m.current = m.current.Parent
	m.rebuildItems()
	m.tableView.SetCursor(m.current.LastCursor)
}

// rebuildItems recomputes the item list from the current directory and
// refreshes the table rows.
func (m *model) rebuildItems() {
	m.items = stexmodel.ComputeItems(m.current, m.cfg)
	m.rebuildRows()
}

// updateInfo refreshes the cached right-pane content for the currently
// selected item. It returns immediately when the selection has not changed.
func (m *model) updateInfo() {
	if !m.ready {
		return
	}
	idx := m.tableView.Cursor()
	if idx < 0 || idx >= len(m.items) {
		return
	}
	item := m.items[idx]
	path := item.FullPath()
	if path == m.infoPath {
		return
	}
	m.infoPath = path

	if item.Kind == stexmodel.TKUpLink {
		m.infoContent = ""
		return
	}

	info := stexmodel.NewFileInfo(path)
	if item.Kind == stexmodel.TKDir {
		children := stexmodel.NewChildrenInfo(item.Dir, maxChildren)
		m.infoContent = stexview.RenderDirInfo(item.Dir, info, children, infoPanelWidth, infoPanelHeight)
	} else {
		m.infoContent = stexview.RenderFileInfo(info, infoPanelWidth, infoPanelHeight)
	}
}

// handleMouseClick selects the clicked row or triggers header sorting.
func (m *model) handleMouseClick(msg tea.MouseClickMsg) (tea.Model, tea.Cmd) {
	if !m.ready {
		return m, nil
	}
	mouse := msg.Mouse()
	clickY := mouse.Y - 1
	clickX := mouse.X - 1

	// In split view, ignore clicks in the right pane
	if m.width-2 >= splitViewThreshold {
		leftWidth := (m.width-2)/2 - 1
		if clickX >= leftWidth+1 {
			return m, nil
		}
	}

	if clickY == 2 {
		m.handleHeaderClick(clickX)
		return m, nil
	}

	rowIndex := clickY - 3
	if rowIndex < 0 || rowIndex >= len(m.items) {
		return m, nil
	}
	m.tableView.SetCursor(rowIndex)
	m.enterSelected()
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
	if !m.ready {
		return m
	}
	if msg.Button == tea.MouseWheelUp {
		m.tableView.MoveUp(scrollLines)
	} else {
		m.tableView.MoveDown(scrollLines)
	}
	return m
}

// View renders the current screen: scan progress or file listing with help
// and optionally a right info panel on wide terminals.
func (m *model) View() tea.View {
	if !m.ready {
		return m.loadingView()
	}

	innerWidth := m.width - 2
	innerHeight := m.height - 2

	if innerWidth >= splitViewThreshold {
		return m.splitView(innerWidth, innerHeight)
	}

	helpStr, helpHeight := m.renderHelp(innerWidth)
	contentHeight := innerHeight - 2 - helpHeight
	if contentHeight < 0 {
		contentHeight = 0
	}

	lines := make([]string, innerHeight)
	lines[0] = stexview.Title(m.current, innerWidth, m.cfg.IconStyle, stexmodel.GroupingString(m.cfg.Grouping))
	lines[1] = dimStyle.Render(strings.Repeat("─", innerWidth))

	fillTable(lines, contentHeight, m.tableView.View())

	helpStart := 2 + contentHeight
	for i := 0; i < helpHeight && helpStart+i < innerHeight; i++ {
		lines[helpStart+i] = helpStr[i]
	}
	for i := helpStart + helpHeight; i < innerHeight; i++ {
		lines[i] = ""
	}

	return m.wrapView(strings.Join(lines, "\n"))
}

// renderHelp returns the rendered help text and its line count.
func (m *model) renderHelp(width int) ([]string, int) {
	content := m.help.View(m.keys)
	if content != "" {
		content = lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(content)
	}
	lines := strings.Split(content, "\n")
	return lines, len(lines)
}

// fillTable copies table lines into the content area of lines.
func fillTable(lines []string, contentHeight int, tableContent string) {
	tableLines := strings.Split(tableContent, "\n")
	for i := 0; i < contentHeight && i < len(tableLines); i++ {
		lines[2+i] = tableLines[i]
	}
	for i := len(tableLines); i < contentHeight; i++ {
		lines[2+i] = ""
	}
}

// wrapView applies the border and screen settings to content.
func (m *model) wrapView(content string) tea.View {
	v := tea.NewView(borderStyle.Render(content))
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

// splitView renders the dual-pane layout with info on the right.
func (m *model) splitView(innerWidth, innerHeight int) tea.View {
	m.updateInfo()

	leftWidth := innerWidth/2 - 1
	rightWidth := innerWidth - leftWidth - 1

	grouping := stexmodel.GroupingString(m.cfg.Grouping)
	titleLine := stexview.Title(m.current, innerWidth, m.cfg.IconStyle, grouping)
	sepLine := dimStyle.Render(strings.Repeat("─", innerWidth))

	helpStr, helpHeight := m.renderHelp(innerWidth)
	contentHeight := innerHeight - 2 - helpHeight
	if contentHeight < 0 {
		contentHeight = 0
	}

	leftContent := renderPadded(m.tableView.View(), contentHeight)

	var rightContent string
	if m.infoContent != "" {
		rightContent = renderPadded(m.infoContent, contentHeight)
	}

	combined := stexview.Split(leftContent, rightContent, leftWidth, rightWidth, dimStyle.Render("│"))

	body := titleLine + "\n" + sepLine + "\n" + combined + "\n" + strings.Join(helpStr, "\n")
	return m.wrapView(body)
}

// renderPadded pads content to exactly height lines.
func renderPadded(content string, height int) string {
	lines := strings.Split(content, "\n")
	out := make([]string, height)
	for i := 0; i < height && i < len(lines); i++ {
		out[i] = lines[i]
	}
	for i := len(lines); i < height; i++ {
		out[i] = ""
	}
	return strings.Join(out, "\n")
}

// loadingView renders the scan progress dialog centred.
func (m *model) loadingView() tea.View {
	innerWidth := m.width - 2
	innerHeight := m.height - loadingMargin
	if innerHeight < minDialogHeight {
		innerHeight = minDialogHeight
	}
	progress := stexview.Progress(m.scanState, innerWidth, innerHeight)
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


