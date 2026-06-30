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
	Sort, Group, Order    key.Binding
	Icons, Hidden         key.Binding
	Search, ClearFilter   key.Binding
}

// explorerKeys is the singleton key map used by the explorer. The list is split into four rows
// for the expanded help view.
var explorerKeys = keys{
	Keys: core.DefaultKeys(),

	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", "l", "right"),
		key.WithHelp("enter/l/→", "open directory"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace", "left"),
		key.WithHelp("esc/←", "go up"),
	),
	Sort: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "sort by name/size"),
	),
	Group: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("g", "grouping mode"),
	),
	Order: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "ascending/descending"),
	),
	Icons: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "toggle icons"),
	),
	Hidden: key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "toggle hidden files"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter by name"),
	),
	ClearFilter: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "clear filter"),
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
		{k.Sort, k.Group, k.Order, k.Icons},
		{k.Hidden, k.Search, k.ClearFilter},
		{k.Help, k.Quit},
	}
}
