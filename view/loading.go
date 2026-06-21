package view

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/model"
)

const (
	progressOverhead   = 9
	commaGroupSize     = 3
	shortenMinLen      = 5
	shortenEllipsisLen = 3
)

// Progress renders the scan progress dialog with current path, item count,
// total size, and any warnings. height is the available lines inside the dialog.
func Progress(state *model.ScanState, width, height int) string {
	state.Mu.Lock()
	currentPath := state.CurrentPath
	total := state.TotalItems
	totalSize := state.TotalSize
	warnings := state.Warnings
	state.Mu.Unlock()

	totalWarnings := len(warnings)

	var body strings.Builder
	body.WriteString(" Scanning...")
	body.WriteString("\n\n")

	itemsStr := commaFormat(total)
	body.WriteString(fmt.Sprintf(" Total items: %s    size: %s", itemsStr, totalSize))
	body.WriteString("\n")

	pathStr := currentPath
	maxPath := width - 4
	if len(pathStr) > maxPath {
		pathStr = shortenPath(pathStr, maxPath)
	}
	body.WriteString(fmt.Sprintf("  %s", pathStr))
	body.WriteString("\n")

	if totalWarnings > 0 {
		body.WriteString("\n")
		avail := height - progressOverhead
		if avail < 1 {
			avail = 1
		}
		body.WriteString(fmt.Sprintf(" WARNING: %d %s - sizes may be inaccurate", totalWarnings, plural("error", totalWarnings)))
		body.WriteString("\n")

		start := 0
		if totalWarnings > avail {
			start = totalWarnings - avail + 1
			body.WriteString(fmt.Sprintf(" ... (%d more) ...", totalWarnings-avail))
			body.WriteString("\n")
		}
		for i := start; i < totalWarnings; i++ {
			wl := warnings[i]
			maxW := width - 2
			if len(wl) > maxW {
				wl = wl[:maxW]
			}
			body.WriteString(fmt.Sprintf(" %s", wl))
			body.WriteString("\n")
		}
		body.WriteString("\n")
	} else {
		body.WriteString("\n")
	}

	body.WriteString(" Press q to abort")

	lines := strings.Split(body.String(), "\n")
	for i, line := range lines {
		if len(line) < width {
			lines[i] = line + strings.Repeat(" ", width-len(line))
		}
	}
	return strings.Join(lines, "\n")
}

// plural returns s when n == 1, otherwise appends "s".
func plural(s string, n int) string {
	if n == 1 {
		return s
	}
	return s + "s"
}

// commaFormat formats n with thousands separators.
func commaFormat(n int64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= commaGroupSize {
		return s
	}
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%commaGroupSize == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return string(result)
}

// shortenPath replaces the middle of a long path with "...".
func shortenPath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	if maxLen < shortenMinLen {
		return path[:maxLen]
	}
	keep := (maxLen - shortenEllipsisLen) / 2
	return path[:keep] + "..." + path[len(path)-keep:]
}
