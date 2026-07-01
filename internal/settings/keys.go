package settings

import "charm.land/bubbles/v2/key"

type keys struct {
	Up, Down, Tab, Confirm, Save, Reset, Cancel key.Binding
}

var settingsKeys = keys{
	Up: key.NewBinding(
		key.WithKeys("up", "k", "w"),
		key.WithHelp("↑/k/w", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j", "s"),
		key.WithHelp("↓/j/s", "down"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "toggle"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "close"),
	),
	Save: key.NewBinding(
		key.WithKeys("S"),
		key.WithHelp("S", "save defaults"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "reset to defaults"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}

func (k keys) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Tab, k.Confirm, k.Save, k.Cancel}
}

func (k keys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Tab},
		{k.Confirm, k.Save, k.Reset, k.Cancel},
	}
}
