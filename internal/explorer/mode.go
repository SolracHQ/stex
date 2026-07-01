// Package explorer implements the main mode: a navigable table of the current directory's
// children with sort, group, hidden, and filter toggles, an info pane for the selected item,
// and a transition to the filter mode.
package explorer

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"github.com/SolracHQ/stex/internal/choose"
	"github.com/SolracHQ/stex/internal/command"
	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/filter"
	"github.com/SolracHQ/stex/internal/settings"

	tea "charm.land/bubbletea/v2"
)

// Mouse click geometry. The border is one cell wide, the header is the second visible line,
// and content starts at the third.
const (
	mouseClickBorderOffset = 1
	mouseClickHeaderRow    = 2
	mouseClickContentStart = 3
	scrollLines            = 3
)

// Explorer is the main mode. Zero value is valid, all state lives in the shared Context.
type Explorer struct{}

// Init resizes the table to the current terminal dimensions, recomputes the item list, and
// rebuilds the table rows.
func (Explorer) Init(ctx *core.Context) tea.Cmd {
	core.Rebuild(ctx)
	return nil
}

// Update handles window resize, key presses, mouse clicks, and mouse wheel. The base view is
// drawn by core.RenderBase, the explorer draws no overlay of its own.
func (Explorer) Update(ctx *core.Context, msg tea.Msg) (core.Mode, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		core.Rebuild(ctx)
		return nil, nil

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, ctx.Keys.Help):
			ctx.Help.ShowAll = !ctx.Help.ShowAll

		case key.Matches(msg, explorerKeys.Up):
			ctx.Table.MoveUp(1)
			core.UpdateInfo(ctx)

		case key.Matches(msg, explorerKeys.Down):
			ctx.Table.MoveDown(1)
			core.UpdateInfo(ctx)

		case key.Matches(msg, explorerKeys.Enter):
			enterSelected(ctx)
			core.UpdateInfo(ctx)

		case key.Matches(msg, explorerKeys.Back):
			goToParent(ctx)
			core.UpdateInfo(ctx)

		case key.Matches(msg, explorerKeys.Group):
			return choose.NewGroupPicker(ctx.Config.Grouping, Explorer{}), nil

		case key.Matches(msg, explorerKeys.Search):
			return filter.New(Explorer{}), nil

		case key.Matches(msg, explorerKeys.Command):
			return command.New(Explorer{}), nil

		case key.Matches(msg, explorerKeys.Settings):
			return settings.New(Explorer{}), nil

		case key.Matches(msg, explorerKeys.ClearFilter):
			if ctx.Config.Filter != nil {
				ctx.Config.Filter = nil
				core.Rebuild(ctx)
			}
		}
		return nil, nil

	case tea.MouseClickMsg:
		handleMouseClick(ctx, msg)
		return nil, nil

	case tea.MouseWheelMsg:
		handleMouseWheel(ctx, msg)
		return nil, nil
	}

	return nil, nil
}

// View returns "". The base view is drawn by core.RenderBase. The explorer itself does not
// draw any overlay, the base renderer handles layout and any active sub mode draws its own
// overlay on top.
func (Explorer) Overlay(_ *core.Context) string { return "" }

// Help returns the explorer's key bindings for the help footer.
func (Explorer) Help() help.KeyMap {
	return explorerKeys
}

// handleMouseClick translates a mouse click into a table action. Clicks on the header row sort
// the table by the clicked column. Clicks on the content rows select the row and navigate
// into directories. Clicks in the right pane of a wide terminal are ignored so the user can
// use the terminal scrollback.
func handleMouseClick(ctx *core.Context, msg tea.MouseClickMsg) {
	if ctx.Current == nil {
		return
	}
	mouse := msg.Mouse()
	clickY := mouse.Y - mouseClickBorderOffset
	clickX := mouse.X - mouseClickBorderOffset

	if ctx.Width-2 >= core.SplitViewThreshold {
		leftWidth := (ctx.Width - 2) / 2
		if clickX >= leftWidth {
			return
		}
	}

	if clickY == mouseClickHeaderRow {
		handleHeaderClick(ctx, clickX)
		return
	}

	rowIndex := clickY - mouseClickContentStart
	if rowIndex < 0 || rowIndex >= len(ctx.Items) {
		return
	}
	ctx.Table.SetCursor(rowIndex)
	enterSelected(ctx)
	core.UpdateInfo(ctx)
}

// handleHeaderClick sorts by the clicked column header. Clicking the active column toggles
// the sort order, clicking a different column switches to sorting by that column.
func handleHeaderClick(ctx *core.Context, clickX int) {
	cols := ctx.Table.Columns()
	x := 1
	for i, col := range cols {
		if clickX >= x && clickX < x+col.Width {
			switch i {
			case 1:
				ctx.Config.SortBy.Toggle()
			case 2:
				if ctx.Config.SortBy == config.SortByName {
					ctx.Config.SortOrder.Toggle()
				} else {
					ctx.Config.SortBy = config.SortByName
				}
			}
			core.Rebuild(ctx)
			return
		}
		x += col.Width
	}
}

// handleMouseWheel scrolls the table by three lines in the wheel direction. core.UpdateInfo is
// called so the right pane follows the new selection.
func handleMouseWheel(ctx *core.Context, msg tea.MouseWheelMsg) {
	if msg.Button == tea.MouseWheelUp {
		ctx.Table.MoveUp(scrollLines)
	} else {
		ctx.Table.MoveDown(scrollLines)
	}
	core.UpdateInfo(ctx)
}
