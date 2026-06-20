package view

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/model"
)

// Progress renders the scan progress dialog with current path, item count,
// total size, and any warnings.
func Progress(state *model.ScanState, width int) string {
	state.Mu.Lock()
	currentPath := state.CurrentPath
	total := state.TotalItems
	totalSize := state.TotalSize
	warnings := state.Warnings
	state.Mu.Unlock()

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
	body.WriteString("\n\n")

	for _, wl := range warnings {
		maxW := width - 2
		if len(wl) > maxW {
			wl = wl[:maxW]
		}
		body.WriteString(fmt.Sprintf(" %s", wl))
		body.WriteString("\n")
	}
	if len(warnings) > 0 {
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

// commaFormat formats n with thousands separators.
func commaFormat(n int64) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}
	var result []byte
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
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
	if maxLen < 5 {
		return path[:maxLen]
	}
	keep := (maxLen - 3) / 2
	return path[:keep] + "..." + path[len(path)-keep:]
}
