package tui

import "charm.land/bubbles/v2/key"

// keyMap defines every key binding in the application and implements
// key.Map so the bubbles/help model can render them.
type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Back  key.Binding
	Sort  key.Binding
	Group key.Binding
	Order key.Binding
	Icons key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// ShortHelp returns the bindings shown in the collapsed help footer.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns all bindings organised into columns for the expanded
// help overlay.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Back},
		{k.Sort, k.Group, k.Order, k.Icons},
		{k.Help, k.Quit},
	}
}

// appKeys is the singleton key map used by the application.
var appKeys = keyMap{
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
		key.WithKeys("esc", "backspace", "h", "left"),
		key.WithHelp("esc/h/←", "go up"),
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
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
