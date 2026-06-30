package core

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
)

// Mode is one state in the app's state machine. A mode owns its key handling and its overlay
// view. The base view is drawn once per frame by the app, the active mode draws an overlay
// that composites on top.
//
// A mode signals a transition by returning a non nil Mode from Update, the app installs it as
// the new active mode and runs its Init. Returning nil means stay in the current mode.
//
// A mode asks the app to do something on its behalf by returning a tea.Cmd that emits one of
// the message types in this package, the app intercepts the message in its Update loop and
// runs the side effect. This keeps long running work in the app and real time behavior in the
// mode.
type Mode interface {
	Init(ctx *Context) tea.Cmd
	Update(ctx *Context, msg tea.Msg) (Mode, tea.Cmd)
	View(ctx *Context) string
	Help() help.KeyMap
}
