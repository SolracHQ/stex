package view

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/model"
)

// Title renders the title bar in the format:
//
//	mode | 1.23 GB | 42f - 3d | /path...
//
// The grouping mode is shown first. Size is never truncated. The path is
// shortened from the left with /.../ when the line exceeds width.
func Title(dir *model.Dir, width int, iconStyle model.IconStyle, grouping string) string {
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

	// Build left to right: mode | size | counts | path
	modePart := grouping
	rightPart := sizeStr + sep + countsStr // protected part: size + counts

	// Check if mode + size + counts fits
	prefix := modePart + sep + rightPart
	if len(prefix)+len(sep) >= width {
		// Narrow: show just mode | size
		pathPart := dir.FullPath
		avail := width - len(modePart) - len(sep) - len(sizeStr) - len(sep)
		if avail < 0 {
			avail = 0
		}
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

	// Full: mode | size | counts | path
	pathPart := dir.FullPath
	avail := width - len(prefix) - len(sep)
	if avail < 0 {
		avail = 0
	}
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
