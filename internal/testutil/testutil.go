// Package testutil holds shared test types used across internal packages. The types are
// compiled into the binary but are tiny and have no side effects.
package testutil

import (
	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"github.com/SolracHQ/stex/internal/core"
)

// StubMode is a no op mode that implements core.Mode. Every method returns nil or zero. Use
// as a sentinel return target in tests where the mode transition is not the focus.
type StubMode struct{}

func (StubMode) Init(*core.Context) tea.Cmd                         { return nil }
func (StubMode) Update(*core.Context, tea.Msg) (core.Mode, tea.Cmd) { return nil, nil }
func (StubMode) Overlay(*core.Context) string                 { return "" }
func (StubMode) Help() help.KeyMap                                  { return nil }
