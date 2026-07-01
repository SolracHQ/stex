package explorer

import (
	"charm.land/bubbles/v2/key"
	"github.com/SolracHQ/stex/internal/core"
)

// keys holds the explorer's key bindings plus the embedded global keys. The FullHelp grouping
// (navigation, toggles, search, help) is used by the help footer.
type keys struct {
	core.Keys

	Up, Down, Enter, Back key.Binding
	Group                 key.Binding
	Search, ClearFilter   key.Binding
	Command, Settings     key.Binding
}

// explorerKeys is the singleton key map used by the explorer. The list is split into four rows
// for the expanded help view.
var explorerKeys = keys{
	Keys: core.DefaultKeys(),

	Up: key.NewBinding(
		key.WithKeys("up", "k", "w"),
		key.WithHelp("↑/k/w", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j", "s"),
		key.WithHelp("↓/j/s", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", "l", "right", "d"),
		key.WithHelp("→/l/d", "move in"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace", "left", "h", "a"),
		key.WithHelp("←/h/a", "move out"),
	),
	Group: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "grouping"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter by name"),
	),
	ClearFilter: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "clear filter"),
	),
	Command: key.NewBinding(
		key.WithKeys(":"),
		key.WithHelp(":", "command"),
	),
	Settings: key.NewBinding(
		key.WithKeys("S"),
		key.WithHelp("S", "settings"),
	),
}

// ShortHelp returns the bindings shown in the collapsed footer.
func (k keys) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Back, k.Help, k.Quit}
}

// FullHelp returns the bindings grouped into rows for the expanded help view.
func (k keys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Back},
		{k.Group, k.Settings},
		{k.Search, k.ClearFilter, k.Command},
		{k.Help, k.Quit},
	}
}
