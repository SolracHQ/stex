# stex

Interactive `du` alternative. TUI file tree browser with live scan progress, sorting, grouping, and keyboard-driven navigation.

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
|---|---|
| `↑`/`k` | move up |
| `↓`/`j` | move down |
| `enter` | open directory |
| `esc`/`bksp` | go to parent |
| `s` | toggle sort (name/size) |
| `o` | toggle order (asc/desc) |
| `g` | toggle grouping |
| `?`/`h` | toggle help overlay |
| `q`/`ctrl+c` | quit |

## Build

```shell
go build -o stex .
```
