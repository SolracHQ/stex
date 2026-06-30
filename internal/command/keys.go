package command

import "charm.land/bubbles/v2/key"

type keys struct {
	Confirm key.Binding
	Cancel  key.Binding
}

var commandKeys = keys{
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "run"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
}
