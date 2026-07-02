package core

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/model"
)

// Title renders the title bar that sits above the file listing. The format is, in wide mode,
// "grouping | total size | file count, dir count | current path". In narrow mode the counts are
// dropped, leaving "grouping | total size | current path".
//
// width is the full width of the title bar in cells. The path segment is shortened from the left
// with "/.../" when it does not fit. grouping is a human readable label like "mixed" or "files
// first" that appears on the left side.
func Title(dir *model.Dir, width int, showIcons bool, grouping string) string {
	if dir == nil {
		return ""
	}
	count := dir.Count()
	sizeStr := dir.Size().String()

	var countsStr string
	if showIcons {
		countsStr = fmt.Sprintf("%d 📄 - %d 📁", count.Files, count.Dirs)
	} else {
		countsStr = fmt.Sprintf("%d f - %d d", count.Files, count.Dirs)
	}

	sep := " | "

	modePart := grouping
	rightPart := sizeStr + sep + countsStr

	prefix := modePart + sep + rightPart
	if len(prefix)+len(sep) >= width {
		pathPart := dir.FullPath()
		avail := max(width-len(modePart)-len(sep)-len(sizeStr)-len(sep), 0)
		if len(pathPart) > avail {
			pathPart = shortenTitlePath(pathPart, avail)
		}
		text := modePart + sep + sizeStr + sep + pathPart
		if len(text) < width {
			pad := width - len(text)
			leftPad := pad / 2
			text = strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", pad-leftPad)
		}
		return text
	}

	pathPart := dir.FullPath()
	avail := max(width-len(prefix)-len(sep), 0)
	if len(pathPart) > avail {
		pathPart = shortenTitlePath(pathPart, avail)
	}
	text := prefix + sep + pathPart
	if len(text) < width {
		pad := width - len(text)
		leftPad := pad / 2
		text = strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", pad-leftPad)
	}
	return text
}

// TitleGroup returns the grouping label, with a "(h)" suffix when hidden files are not shown.
// The result is meant to be passed to Title as the grouping argument.
func TitleGroup(ctx *Context) string {
	g := config.GroupingString(ctx.Config.Grouping)
	if !ctx.Config.ShowHidden {
		g += " (h)"
	}
	return g
}

// shortenTitlePath replaces the left side of a long path with "/.../" so the most specific
// part of the path stays visible. Returns "" when maxLen is less than 1.
func shortenTitlePath(path string, maxLen int) string {
	if maxLen < 1 {
		return ""
	}
	if len(path) <= maxLen {
		return path
	}
	keep := max(maxLen-4, 1)
	tail := path[len(path)-keep:]
	if len(tail) > 0 && tail[0] == '/' {
		tail = tail[1:]
	}
	return "/.../" + tail
}
