# stex

stex (Storage Explorer) is a TUI app to explore disk usage. It is heavily
inspired by ncdu and nvim (in the mode architecture). It aims to be simple to
use, visually pleasant and multiplatform.

## Install

```shell
go install github.com/SolracHQ/stex@latest
```

Or clone and build:

```shell
git clone https://github.com/SolracHQ/stex.git
cd stex
go build -o stex .
```

## Usage

```shell
stex [flags] [path]
```

Defaults to the current directory when no path is given.

### Flags

| Flag | Short | Description |
| --- | --- | --- |
| `--icons` | `-i` | start with emoji icons enabled |
| `--show-all` | `-a` | start with hidden files shown |
| `--no-live-filter` | `-L` | disable live filter (compile on enter) |
| `--help` | `-h` | show this help |

## What makes it different?

stex is built around the concept of modes. Each mode changes the app behavior,
the keybindings and the overlay, reducing the amount of keybindings per mode.

### Modes

#### Explorer Mode

Explorer mode is the main app mode. It allows you to navigate the filesystem
subtree and select items to act over. It is focused completely on movement,
allowing three kinds of movement sets, hjkl for vim users, arrows for
traditionalists and awsd for gamers.

| Key | Action |
| --- | --- |
| `↑`/`k`/`w` | move up |
| `↓`/`j`/`s` | move down |
| `→`/`l`/`d` | open directory |
| `←`/`h`/`a` | go to parent |
| `/` | filter mode |
| `:` | command mode |
| `S` | settings panel |
| `g` | grouping picker |
| `c` | clear filter |
| `?` | toggle help |
| `ctrl+c` | quit |

#### Filter Mode

Filter mode lets you search the current directory by regex. Press `/` to open a
search bar at the bottom, all keystrokes are directed to the input and normal
navigation is suspended. The list narrows in real time as you type. On slower
machines you can disable live filter with the `--no-live-filter` flag or toggle
it at runtime with `ctrl+l`, in manual mode the filter only compiles on enter.

| Key | Action |
| --- | --- |
| type a character | append to the pattern and filter live |
| `backspace` | remove the last character |
| `enter` | keep the filter and return to explorer |
| `esc` | clear the filter and return to explorer |
| `ctrl+l` | toggle live filter |

#### Settings Panel

Settings panel is the UI to customize the app behavior. You can toggle sort,
order, grouping, icons, hidden files and live filter. It is the same
configuration you can do via CLI flags or command mode, just in a friendly
menu. There are plans to add UI theme customization in the future.

| Key | Action |
| --- | --- |
| `↑`/`k`/`w` | move up |
| `↓`/`j`/`s` | move down |
| `tab` | toggle focused row |
| `enter` | close and keep changes |
| `r` | reset to defaults |
| `S` | save defaults to config file |
| `esc` | revert and cancel |

#### Command Mode

Command mode is the command centric way to configure the app, allowing
everything the settings panel does but purely through commands, like a command
palette in any code editor. The idea is for it to grow into a more powerful
tool to navigate and act on the filesystem without leaving the keyboard. See
the full command list in [docs/commands.md](docs/commands.md). I even borrowed
`:q` to quit, old habits die hard.

| Key | Action |
| --- | --- |
| type | write the command |
| `tab` | accept suggestion |
| `enter` | run the command |
| `↑`/`k`/`w` | previous suggestion |
| `↓`/`j`/`s` | next suggestion |
| `esc` | cancel |

## Configuration

stex reads config from `$XDG_CONFIG/stex/config.json` when it exists but CLI
flags take precedence over the file and the file takes precedence over
defaults. That way you can set your preferred defaults once in the config file
and override per context with shell aliases like `alias sxh='stex -a'` for
checking disk usage in your home directory where `.local` is usually the
culprit.

## Architecture

stex uses a mode architecture borrowed from nvim. The app owns a shared
Context that holds the file tree, configuration and the table widget. Each
mode defines its own keybindings and draws an optional overlay on top of the
base view. The base view (table, info pane, title, footer) is drawn by the app
once per frame, the active mode composites its overlay on top. This keeps the
explorer lean and makes adding new features a matter of writing a new mode
package.

## License

MIT. See [LICENSE](LICENSE).
