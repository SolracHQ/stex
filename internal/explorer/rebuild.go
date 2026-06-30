package explorer

import (
	"github.com/SolracHQ/stex/internal/core"
	stexmodel "github.com/SolracHQ/stex/internal/model"
)

// Layout constants for the right pane. Width and height are the cell budget for the rendered
// info, maxTopChildren bounds the "Largest Children" segment in the directory info.
const (
	infoPanelWidth  = 80
	infoPanelHeight = 10
	maxTopChildren  = 10
)

// Rebuild resizes the table to the current terminal dimensions, recomputes the item list from
// the current directory and the config, and rebuilds the table columns and rows. Call after
// any state change that affects the listing, sort, group, filter, hidden toggle, navigation
// into a directory.
func Rebuild(ctx *core.Context) {
	applyResize(ctx)
	items := ComputeItems(ctx.Current, ctx.Config)
	ctx.Items = items
	buildColumns(ctx)
	buildRows(ctx, items)
}

// applyResize fits the table to the current terminal dimensions. Safe to call when dims are
// zero, it does nothing in that case so the program does not crash before the first
// WindowSizeMsg.
func applyResize(ctx *core.Context) {
	if ctx.Width == 0 || ctx.Height == 0 {
		return
	}
	innerWidth := ctx.Width - 2
	innerHeight := ctx.Height - 2
	ctx.Table.SetWidth(innerWidth)
	ctx.Table.SetHeight(max(innerHeight-3, 1))
}

// updateInfo refreshes the cached right pane content for the current cursor position. Returns
// immediately when the selection has not changed, so wheel scrolling stays cheap.
func updateInfo(ctx *core.Context) {
	if ctx.Current == nil {
		return
	}
	idx := ctx.Table.Cursor()
	if idx < 0 || idx >= len(ctx.Items) {
		return
	}
	item := ctx.Items[idx]
	path := item.FullPath()
	if path == ctx.Info.Path {
		return
	}
	ctx.Info.Path = path

	if _, ok := item.(*stexmodel.UpLink); ok {
		ctx.Info.Content = ""
		return
	}

	info := core.NewFileInfo(path)
	if dir, ok := item.(*stexmodel.Dir); ok {
		children := TopChildren(dir, maxTopChildren)
		ctx.Info.Content = core.RenderDirInfo(info, dir.Size(), children, infoPanelWidth, infoPanelHeight)
	} else {
		ctx.Info.Content = core.RenderFileInfo(info, infoPanelWidth, infoPanelHeight)
	}
}
