package core

// Rebuild resizes the table to the current terminal dimensions, recomputes the item list from
// the current directory and the config, and rebuilds the table columns and rows. Call after
// any state change that affects the listing, sort, group, filter, hidden toggle, navigation
// into a directory.
func Rebuild(ctx *Context) {
	applyResize(ctx)
	items := ctx.Current.ComputeItems(ctx.Config)
	ctx.Items = items
	buildColumns(ctx)
	buildRows(ctx, items)
}

// applyResize fits the table to the current terminal dimensions. Safe to call when dims are
// zero, it does nothing in that case so the program does not crash before the first
// WindowSizeMsg.
func applyResize(ctx *Context) {
	if ctx.Width == 0 || ctx.Height == 0 {
		return
	}
	innerWidth := ctx.Width - 2
	innerHeight := ctx.Height - 2
	ctx.Table.SetWidth(innerWidth)
	ctx.Table.SetHeight(max(innerHeight-3, 1))
}
