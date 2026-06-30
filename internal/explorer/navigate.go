package explorer

import (
	"github.com/SolracHQ/stex/internal/core"
	stexmodel "github.com/SolracHQ/stex/internal/model"
)

// enterSelected opens the currently selected directory or navigates up when the cursor is on
// the ".." entry. Files are a no op. The active filter is cleared on every navigation so the
// user lands in a new directory without the old filter hiding its children.
func enterSelected(ctx *core.Context) {
	ctx.Config.Filter = nil
	ctx.Info.Path = ""
	if len(ctx.Items) == 0 {
		return
	}
	idx := ctx.Table.Cursor()
	if idx < 0 || idx >= len(ctx.Items) {
		return
	}
	ctx.Current.SetLastSelectedUID(ctx.Items[idx].UID())

	switch item := ctx.Items[idx].(type) {
	case *stexmodel.Dir:
		ctx.Current = item
	case *stexmodel.UpLink:
		ctx.Current = item.ParentDir()
	default:
		return
	}
	Rebuild(ctx)
	restoreCursorByUID(ctx)
}

// goToParent moves the view up one directory level. No op at the root. Like enterSelected, the
// filter is cleared so the destination directory is shown in full.
func goToParent(ctx *core.Context) {
	ctx.Config.Filter = nil
	ctx.Info.Path = ""
	if ctx.Current.ParentDir() == nil {
		return
	}
	idx := ctx.Table.Cursor()
	if idx >= 0 && idx < len(ctx.Items) {
		ctx.Current.SetLastSelectedUID(ctx.Items[idx].UID())
	}
	ctx.Current = ctx.Current.ParentDir()
	Rebuild(ctx)
	restoreCursorByUID(ctx)
}

// restoreCursorByUID positions the cursor on the row whose UID matches the current directory's
// lastSelectedUID, when one is set. The set call happens just before navigation in
// enterSelected and goToParent, so the cursor lands back on the same logical row the user came
// from.
func restoreCursorByUID(ctx *core.Context) {
	uid := ctx.Current.LastSelectedUID()
	if uid == 0 {
		return
	}
	for i, item := range ctx.Items {
		if item.UID() == uid {
			ctx.Table.SetCursor(i)
			return
		}
	}
}
