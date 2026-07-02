// Package choose implements a generic options dialog. The caller supplies a title and a list
// of Option values. Each Option's Action is invoked with the shared Context when the user
// confirms the selection and returns the next Mode to transition to.
package choose

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"

	tea "charm.land/bubbletea/v2"
)

// Option is a single selectable item in the dialog. Label is shown in the list, Action is
// called when the user confirms this option.
type Option struct {
	Label  string
	Action func(ctx *core.Context) (core.Mode, tea.Cmd)
}

// Choose is a generic options dialog with cursor navigation.
type Choose struct {
	title   string
	options []Option
	cursor  int
	backTo  core.Mode
}

// New returns a Choose dialog with the given title and options. backTo is the mode to return
// to on cancel.
func New(title string, options []Option, backTo core.Mode) *Choose {
	return &Choose{title: title, options: options, backTo: backTo}
}

// SetCursor positions the highlight on the given row. Values outside the option range are
// silently ignored.
func (ch *Choose) SetCursor(i int) {
	if i >= 0 && i < len(ch.options) {
		ch.cursor = i
	}
}

// Init returns nil.
func (ch *Choose) Init(_ *core.Context) tea.Cmd { return nil }

// Update moves the cursor, confirms the selection, or cancels the dialog.
func (ch *Choose) Update(ctx *core.Context, msg tea.Msg) (core.Mode, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, chooseKeys.Up):
			if ch.cursor > 0 {
				ch.cursor--
			}
		case key.Matches(msg, chooseKeys.Down):
			if ch.cursor < len(ch.options)-1 {
				ch.cursor++
			}
		case key.Matches(msg, chooseKeys.Confirm):
			if ch.cursor >= 0 && ch.cursor < len(ch.options) {
				return ch.options[ch.cursor].Action(ctx)
			}
		case key.Matches(msg, chooseKeys.Cancel):
			if ch.backTo != nil {
				return ch.backTo, nil
			}
			return nil, nil
		}
	}
	return nil, nil
}

// Help returns the choose key bindings for the help footer.
func (ch *Choose) Help() help.KeyMap {
	return core.FlatKeyMap{chooseKeys.Up, chooseKeys.Down, chooseKeys.Confirm, chooseKeys.Cancel}
}

// NewGroupPicker returns a Choose dialog pre configured with the five grouping values (Mixed,
// FilesFirst, DirsFirst, FilesOnly, DirsOnly). The cursor is pre set on current. backTo is
// the mode to return to on cancel. On confirm the picked value is written to ctx.Config.Grouping.
func NewGroupPicker(current config.Grouping, backTo core.Mode) *Choose {
	values := []config.Grouping{config.Mixed, config.FilesFirst, config.DirsFirst, config.FilesOnly, config.DirsOnly}
	options := make([]Option, len(values))
	for i, v := range values {
		v := v
		options[i] = Option{
			Label: config.GroupingString(v),
			Action: func(ctx *core.Context) (core.Mode, tea.Cmd) {
				ctx.Config.Grouping = v
				core.Rebuild(ctx)
				return backTo, nil
			},
		}
	}
	ch := New("Group By", options, backTo)
	for i, v := range values {
		if v == current {
			ch.SetCursor(i)
			break
		}
	}
	return ch
}
