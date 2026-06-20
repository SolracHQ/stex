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
stex [path]
```

Defaults to the current directory when no path is given.

### Keys

| Key | Action |
|---|---|---|
| `↑`/`k` | move up |
| `↓`/`j` | move down |
| `enter`/`l`/`right` | open directory |
| `esc`/`h`/`left` | go to parent |
| `s` | toggle sort (name/size) |
| `o` | toggle order (asc/desc) |
| `g` | toggle grouping |
| `i` | toggle emoji icons |
| `?` | toggle help overlay |
| `q`/`ctrl+c` | quit |

## Build

```shell
go build -o stex .
```
