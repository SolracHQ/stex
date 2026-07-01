// Package settings implements the settings panel mode. It owns a cursor over a list of rows
// and a config snapshot taken at Init for revert. Changes apply to ctx.Config in real time.
// The grouping row opens a picker sub dialog.
package settings

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"github.com/SolracHQ/stex/internal/choose"
	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"

	tea "charm.land/bubbletea/v2"
)

const (
	rowSort = iota
	rowOrder
	rowGroup
	rowIcons
	rowHidden
	rowLiveFilter
	rowCount
)

// Settings is the settings panel mode. It owns a cursor and a config snapshot. The snapshot
// is taken at Init so esc can revert. Zero value is not valid, use New.
type Settings struct {
	snapshot config.Config
	cursor   int
	inited   bool
	returnTo core.Mode
}

// New returns a Settings bound to returnTo. The caller installs it as the active mode.
func New(returnTo core.Mode) *Settings {
	return &Settings{returnTo: returnTo}
}

// Init takes a snapshot of the current config for revert on esc. Only runs once.
func (s *Settings) Init(ctx *core.Context) tea.Cmd {
	if !s.inited {
		s.snapshot = ctx.Config
		s.inited = true
	}
	return nil
}

// Update handles the settings key map. Tab toggles the focused row or opens the grouping
// picker, enter closes, esc reverts and closes, S saves and closes, R resets.
func (s *Settings) Update(ctx *core.Context, msg tea.Msg) (core.Mode, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, settingsKeys.Up):
			if s.cursor > 0 {
				s.cursor--
			}
		case key.Matches(msg, settingsKeys.Down):
			if s.cursor < rowCount-1 {
				s.cursor++
			}
		case key.Matches(msg, settingsKeys.Tab):
			if next := s.applyFocused(ctx); next != nil {
				return next, nil
			}
			core.Rebuild(ctx)
			return nil, nil
		case key.Matches(msg, settingsKeys.Confirm):
			return s.returnTo, nil
		case key.Matches(msg, settingsKeys.Save):
			save(ctx.Config)
			return s.returnTo, nil
		case key.Matches(msg, settingsKeys.Reset):
			ctx.Config = config.DefaultConfig()
			core.Rebuild(ctx)
			return nil, nil
		case key.Matches(msg, settingsKeys.Cancel):
			ctx.Config = s.snapshot
			return s.returnTo, nil
		}
	}
	return nil, nil
}

// Help returns the settings key bindings for the help footer.
func (s *Settings) Help() help.KeyMap {
	return settingsKeys
}

// applyFocused either toggles the focused field or opens the grouping picker for the group
// row. Returns the new mode when the grouping picker is opened, nil otherwise.
func (s *Settings) applyFocused(ctx *core.Context) core.Mode {
	switch s.cursor {
	case rowSort:
		ctx.Config.SortBy.Toggle()
	case rowOrder:
		ctx.Config.SortOrder.Toggle()
	case rowGroup:
		return choose.NewGroupPicker(ctx.Config.Grouping, s)
	case rowIcons:
		ctx.Config.ShowIcons = !ctx.Config.ShowIcons
	case rowHidden:
		ctx.Config.ShowHidden = !ctx.Config.ShowHidden
	case rowLiveFilter:
		ctx.Config.LiveFilter = !ctx.Config.LiveFilter
	}
	return nil
}

// save writes the current config to the user's config file.
func save(cfg config.Config) error {
	return cfg.Save()
}
