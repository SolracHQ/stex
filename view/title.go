package view

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/model"
)

// Title renders the title bar in the format:
//
//	1.23 GB | 42f - 3d | /path...
//
// Size is always fully visible on the left. The path is truncated from the
// left with /.../ when the line exceeds width.
func Title(dir *model.Dir, width int, iconStyle model.IconStyle) string {
	count := dir.Count()
	sizeStr := dir.Size().String()

	var countsStr string
	switch iconStyle {
	case model.IconLetters:
		countsStr = fmt.Sprintf("%d f - %d d", count.Files, count.Dirs)
	case model.IconEmoji:
		countsStr = fmt.Sprintf("%d 📄 - %d 📁", count.Files, count.Dirs)
	}

	sep := " | "

	// Build left to right: size | counts | path
	leftPart := sizeStr + sep + countsStr

	if len(leftPart)+len(sep) >= width {
		// Not enough room for counts — show just size
		pathPart := dir.FullPath
		avail := width - len(sizeStr) - len(sep)
		if avail < 0 {
			avail = 0
		}
		if len(pathPart) > avail {
			pathPart = shortenTitlePath(pathPart, avail)
		}
		text := sizeStr + sep + pathPart
		if len(text) < width {
			pad := width - len(text)
			leftPad := pad / 2
			text = strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", pad-leftPad)
		}
		return text
	}

	pathPart := dir.FullPath
	avail := width - len(leftPart) - len(sep)
	if avail < 0 {
		avail = 0
	}
	if len(pathPart) > avail {
		pathPart = shortenTitlePath(pathPart, avail)
	}
	text := leftPart + sep + pathPart
	if len(text) < width {
		pad := width - len(text)
		leftPad := pad / 2
		text = strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", pad-leftPad)
	}
	return text
}

// shortenTitlePath replaces the left side of a long path with /.../.
func shortenTitlePath(path string, maxLen int) string {
	if maxLen < 1 {
		return ""
	}
	if len(path) <= maxLen {
		return path
	}
	keep := maxLen - 4
	if keep < 1 {
		keep = 1
	}
	tail := path[len(path)-keep:]
	if len(tail) > 0 && tail[0] == '/' {
		tail = tail[1:]
	}
	return "/.../" + tail
}
