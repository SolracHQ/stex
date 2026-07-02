// Package command implements the ":" command mode. The user types a line, the parser walks a
// commands map, and most commands mutate ctx.Config and return to the caller.
package command

import (
	"strconv"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"

	tea "charm.land/bubbletea/v2"
)

type cmdRun func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd)

type cmdDef struct {
	args []string
	run  cmdRun
}

var commands = map[string]cmdDef{
	"quit": {
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			return nil, tea.Quit
		},
	},
	"save": {
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			_ = ctx.Config.Save()
			return returnTo, nil
		},
	},
	"sort": {
		args: []string{"ascending", "descending"},
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			switch arg {
			case "asc", "ascending":
				ctx.Config.SortOrder = config.Ascending
			case "desc", "descending":
				ctx.Config.SortOrder = config.Descending
			default:
				return returnTo, nil
			}
			return returnTo, nil
		},
	},
	"sortby": {
		args: []string{"name", "size"},
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			switch arg {
			case "name":
				ctx.Config.SortBy = config.SortByName
			case "size":
				ctx.Config.SortBy = config.SortBySize
			default:
				return returnTo, nil
			}
			return returnTo, nil
		},
	},
	"group": {
		args: []string{"files", "dirs", "filesonly", "dirsonly", "mixed"},
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			switch arg {
			case "files":
				ctx.Config.Grouping = config.FilesFirst
			case "dirs":
				ctx.Config.Grouping = config.DirsFirst
			case "filesonly":
				ctx.Config.Grouping = config.FilesOnly
			case "dirsonly":
				ctx.Config.Grouping = config.DirsOnly
			case "mixed":
				ctx.Config.Grouping = config.Mixed
			default:
				return returnTo, nil
			}
			return returnTo, nil
		},
	},
	"toggle": {
		args: []string{"icons", "hidden", "live"},
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			switch arg {
			case "icons":
				ctx.Config.ShowIcons = !ctx.Config.ShowIcons
			case "hidden":
				ctx.Config.ShowHidden = !ctx.Config.ShowHidden
			case "live":
				ctx.Config.LiveFilter = !ctx.Config.LiveFilter
			default:
				return returnTo, nil
			}
			return returnTo, nil
		},
	},
	"up": {
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			n := 1
			if arg != "" {
				parsed, err := strconv.Atoi(arg)
				if err != nil || parsed < 1 {
					return returnTo, nil
				}
				n = parsed
			}
			for range n {
				if ctx.Current.ParentDir() == nil {
					break
				}
				ctx.Current = ctx.Current.ParentDir()
			}
			core.Rebuild(ctx)
			return returnTo, nil
		},
	},
}

// commandVerbs is the canonical verb list shown when no argument has been typed yet. Short
// forms like q are handled by the parser but never appear in suggestions.
var commandVerbs = []string{"quit", "save", "sort", "sortby", "group", "toggle", "up"}

// Command is the command line mode. It owns a textinput widget with tab completion and a return
// target that the mode transitions back to on complete or cancel.
type Command struct {
	input    textinput.Model
	returnTo core.Mode
}

// New returns a Command bound to returnTo. The caller installs it as the active mode.
func New(returnTo core.Mode) *Command {
	input := textinput.New()
	input.Prompt = ":"
	input.Placeholder = "command"
	input.ShowSuggestions = true
	input.SetSuggestions(commandVerbs)
	input.KeyMap.AcceptSuggestion = key.NewBinding(key.WithKeys("tab"))
	input.KeyMap.NextSuggestion = key.NewBinding(key.WithKeys("down", "j"))
	input.KeyMap.PrevSuggestion = key.NewBinding(key.WithKeys("up", "k"))

	return &Command{input: input, returnTo: returnTo}
}

// Init focuses the textinput so the user can type the command.
func (c *Command) Init(_ *core.Context) tea.Cmd {
	return c.input.Focus()
}

// Update processes the textinput, refreshes the suggestion pool, and handles enter (run
// command) and esc (cancel).
func (c *Command) Update(ctx *core.Context, msg tea.Msg) (core.Mode, tea.Cmd) {
	preValue := c.input.Value()
	needRefresh := false

	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case key.Matches(keyMsg, commandKeys.Confirm):
			if cur := c.input.CurrentSuggestion(); cur != "" {
				c.input.SetValue(cur)
			}
			return runCommand(ctx, c.input.Value(), c.returnTo)
		case key.Matches(keyMsg, commandKeys.Cancel):
			return c.returnTo, nil
		}

		found := strings.Contains(preValue, " ")
		afterSpace := strings.IndexByte(c.input.Value(), ' ')
		if (!found && afterSpace != -1) || (found && afterSpace == -1) {
			needRefresh = true
		}
	}

	var cmd tea.Cmd
	c.input, cmd = c.input.Update(msg)

	if needRefresh {
		c.refreshSuggestions(c.input.Value())
	}

	if _, ok := msg.(tea.KeyPressMsg); !ok {
		c.refreshSuggestions(c.input.Value())
	}

	return nil, cmd
}

// refreshSuggestions sets the suggestion pool based on whether the current value has a space
// (verb plus arguments) or not (verb only).
func (c *Command) refreshSuggestions(value string) {
	before, _, ok := strings.Cut(value, " ")
	if !ok {
		c.input.SetSuggestions(commandVerbs)
		return
	}

	verb := before
	def, exists := commands[verb]
	if !exists || len(def.args) == 0 {
		c.input.SetSuggestions(nil)
		return
	}

	full := make([]string, len(def.args))
	for i, arg := range def.args {
		full[i] = verb + " " + arg
	}
	c.input.SetSuggestions(full)
}

// Help returns the command key bindings for the help footer.
func (c *Command) Help() help.KeyMap {
	return core.FlatKeyMap{commandKeys.Confirm, commandKeys.Cancel}
}

func runCommand(ctx *core.Context, value string, returnTo core.Mode) (core.Mode, tea.Cmd) {
	parts := splitFields(value)
	if len(parts) == 0 {
		return returnTo, nil
	}

	verb, arg := parts[0], ""
	if len(parts) > 1 {
		arg = parts[1]
	}

	def, exists := commands[verb]
	if !exists {
		return returnTo, nil
	}

	if len(def.args) > 0 && arg == "" {
		return returnTo, nil
	}

	return def.run(ctx, arg, returnTo)
}

func splitFields(value string) []string {
	var out []string
	var cur strings.Builder
	for _, r := range value {
		if r == ' ' || r == '\t' {
			if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
			continue
		}
		cur.WriteRune(r)
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}
