# Command Mode Reference

## Verbs

### quit

Quit the app.

### save

Persist the current configuration to `$XDG_CONFIG/stex/config.json`.

### sort

Set the sort order. Requires one argument.

| Argument | Aliases | Effect |
| --- | --- | --- |
| `ascending` | `asc` | smallest first |
| `descending` | `desc` | largest first |

### sortby

Set the sort field. Requires one argument.

| Argument | Effect |
| --- | --- |
| `name` | sort alphabetically |
| `size` | sort by file size |

### group

Set the grouping mode. Requires one argument.

| Argument | Effect |
| --- | --- |
| `files` | files first, then directories |
| `dirs` | directories first, then files |
| `filesonly` | only show files |
| `dirsonly` | only show directories |
| `mixed` | interleave files and directories sorted together |

### toggle

Toggle a boolean setting. Requires one argument.

| Argument | Effect |
| --- | --- |
| `icons` | show or hide emoji icons |
| `hidden` | show or hide hidden files |
| `live` | enable or disable live filter |

### up

Navigate up one or more directories. With no arguments goes to the parent.
With a number N goes up N levels, stopping at the root.
