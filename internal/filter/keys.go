package filter

import "charm.land/bubbles/v2/key"

// keys holds the filter mode's key bindings.
type keys struct {
	Confirm    key.Binding
	Cancel     key.Binding
	ToggleLive key.Binding
}

// filterKeys is the singleton key map used in filter mode. The
// bindings are exported via the FlatKeyMap returned from Help.
var filterKeys = keys{
	Confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "keep filter"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "clear filter"),
	),
	ToggleLive: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "toggle live filter"),
	),
}
