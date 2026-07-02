// Package filter implements the regex filter mode. It owns a textinput widget and writes the
// compiled pattern to ctx.Config.Filter either live (on every keystroke) or on enter (manual
// mode).
package filter

import (
	"regexp"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	"github.com/SolracHQ/stex/internal/core"

	tea "charm.land/bubbletea/v2"
)

// Filter is the regex filter mode. It owns a textinput widget and updates the shared config
// filter on every keystroke when LiveFilter is on, or on Enter when LiveFilter is off. The
// returnTo field is the mode the filter transitions back to on confirm or cancel.
type Filter struct {
	input    textinput.Model
	returnTo core.Mode
}

// New returns a fresh Filter with the textinput focused. returnTo is the mode the filter
// returns to when the user confirms or cancels, the caller is responsible for picking it.
func New(returnTo core.Mode) *Filter {
	input := textinput.New()
	input.Prompt = "/"
	input.Placeholder = "regex"
	return &Filter{input: input, returnTo: returnTo}
}

// Init focuses the textinput so the user can type the pattern.
func (flt *Filter) Init(_ *core.Context) tea.Cmd {
	return flt.input.Focus()
}

// Update handles the filter key map plus the underlying text input. Live mode commits the
// pattern on every keystroke, manual mode waits for enter.
func (flt *Filter) Update(ctx *core.Context, msg tea.Msg) (core.Mode, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, filterKeys.Confirm):
			commit(ctx, flt.input.Value())
			return flt.returnTo, nil
		case key.Matches(msg, filterKeys.Cancel):
			ctx.Config.Filter = nil
			return flt.returnTo, nil
		case key.Matches(msg, filterKeys.ToggleLive):
			ctx.Config.LiveFilter = !ctx.Config.LiveFilter
			if ctx.Config.LiveFilter {
				commit(ctx, flt.input.Value())
			} else {
				ctx.Config.Filter = nil
			}
			core.Rebuild(ctx)
			return nil, nil
		}
	}
	var cmd tea.Cmd
	flt.input, cmd = flt.input.Update(msg)
	if ctx.Config.LiveFilter {
		commit(ctx, flt.input.Value())
		core.Rebuild(ctx)
		return nil, cmd
	}
	return nil, cmd
}

// Help returns the filter key bindings for the help footer.
func (flt *Filter) Help() help.KeyMap {
	return core.FlatKeyMap{filterKeys.Confirm, filterKeys.Cancel, filterKeys.ToggleLive}
}

// commit compiles pattern and stores it on ctx.Config.Filter. An empty pattern clears the
// filter. A pattern that fails to compile is silently dropped, the previous filter is kept.
func commit(ctx *core.Context, pattern string) {
	if pattern == "" {
		ctx.Config.Filter = nil
		return
	}
	if re, err := regexp.Compile(pattern); err == nil {
		ctx.Config.Filter = re
	}
}
