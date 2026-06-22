# stex

Interactive TUI to explore where disk space is being used. Scan any directory, sort by size, follow the color gradient (red = biggest) to find what's consuming storage. Live scan progress and grouping. Rich keyboard navigation for SSH sessions, mouse support for local usage.

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
|---|---|---|
| `--icons` | `-i` | start with emoji icons enabled |
| `--show-all` | `-a` | start with hidden files shown |
| `--no-live-filter` | `-L` | disable live filter (compile on enter) |
| `--help` | `-h` | show this help |

### Modes

stex has two key modes: **explore** (default) and **filter** (triggered by `/`).

#### Explore mode

Navigate, sort, group, and toggle settings.

| Key | Action |
|---|---|
| `↑`/`k` | move up |
| `↓`/`j` | move down |
| `enter`/`l`/`→` | open directory |
| `esc`/`←` | go to parent |
| `s` | toggle sort (name/size) |
| `o` | toggle order (asc/desc) |
| `g` | toggle grouping |
| `i` | toggle emoji icons |
| `H` | toggle hidden files |
| `/` | enter filter mode |
| `c` | clear active filter |
| `?` | toggle help overlay |
| `q`/`ctrl+c` | quit |

#### Filter mode

Press `/` to open a search bar at the bottom. All keystrokes are directed to the filter textbox, normal navigation and commands are suspended. The item list narrows in real time as you type a regex pattern.

| Key | Action |
|---|---|
| type a character | append to the regex pattern and filter live |
| `backspace` | remove the last character |
| `enter` | keep the filter (and compile it in manual mode) and return to explore mode |
| `esc` | clear the filter and return to explore mode |
| `ctrl+l` | toggle live filter (on by default) / manual mode (compile on enter) |

## Build

```shell
go build -o stex .
```
