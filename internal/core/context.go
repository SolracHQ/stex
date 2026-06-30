// Package core is the app's base architecture. The base view is drawn once per frame, modes
// contribute a keymap and can return an overlay that composites on top. The model is the nvim
// style, a stable base, modes that change keys, overlays that add context, modes that compose.
//
// The app is the long lived piece, it owns the Context and any program lifetime orchestration.
// Modes are transient states for real time behavior. A sync that needs to lock input and mouse
// is the exception, it lives in a mode so the lock is enforced by the mode being active.
package core

import (
	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/model"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/table"
)

// Context is the mutable state bag passed to every mode. Modes read and write the fields
// directly. The shared widget models (Table, Help, Spinner) and the scanned tree (Root, Current)
// live here so a mode transition does not need to rebuild them.
type Context struct {
	Width, Height int
	Path          string

	Root, Current *model.Dir
	Config        config.Config

	Table table.Model
	Info  InfoState
	Items []model.FileSystemItem

	Help help.Model
	Keys Keys

	Ready bool
}

// InfoState holds the cached right pane content. Tracking the path lets the explorer skip
// re rendering when the cursor has not moved.
type InfoState struct {
	Path    string
	Content string
}
