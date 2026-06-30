package core

import "charm.land/bubbles/v2/key"

// Keys holds the key bindings that work in every mode. Only bindings that must be active
// regardless of the current mode belong here, the quit and help toggles. Mode specific keys
// live in the mode package and are exposed through the Mode's Help method.
//
// q is intentionally not a global quit so text input modes like filter can receive it without
// the app intercepting it first.
type Keys struct {
	Quit key.Binding
	Help key.Binding
}

// DefaultKeys returns the standard key bindings, ctrl+c and ctrl+d to quit, ? to toggle the
// help footer.
func DefaultKeys() Keys {
	return Keys{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "ctrl+d"),
			key.WithHelp("ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}

// FlatKeyMap wraps a flat slice of bindings into a help.KeyMap. ShortHelp returns every binding
// in declaration order, FullHelp returns them as a single row. Modes that have a small flat
// set of bindings (filter, command, choose) use this directly instead of defining a separate
// type.
type FlatKeyMap []key.Binding

// ShortHelp returns the underlying slice of bindings.
func (km FlatKeyMap) ShortHelp() []key.Binding { return km }

// FullHelp returns the bindings as a single row, so the help footer renders them on one line.
func (km FlatKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{km}
}
