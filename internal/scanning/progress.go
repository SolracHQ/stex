package scanning

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/internal/model"
	"github.com/SolracHQ/stex/internal/styles"
)

// Layout constants for the progress dialog. progressOverhead is the number of lines the
// dialog always reserves for header, totals, and footer, the warning list is then given
// whatever vertical room remains. The two shorten helpers control the middle ellipsis path
// display.
const (
	progressOverhead   = 9
	commaGroupSize     = 3
	shortenMinLen      = 5
	shortenEllipsisLen = 3
	loadingMargin      = 6
	minDialogHeight    = 10
)

// progressBox renders the progress dialog framed by the scan border style. width and height
// are the full terminal dimensions, the dialog is centered inside.
func progressBox(state *model.ScanState, width, height int) string {
	innerWidth := width - 2
	innerHeight := max(height-loadingMargin, minDialogHeight)
	progress := progressBody(state, innerWidth, innerHeight)
	return styles.BorderNorm.Render(progress)
}

// progressBody renders the body of the progress dialog. It shows the current path, the running
// item count, the running total size, and the last few warnings when the scan has encountered
// permission problems.
func progressBody(state *model.ScanState, width, height int) string {
	state.Mu.Lock()
	currentPath := state.CurrentPath
	total := state.TotalItems
	totalSize := state.TotalSize
	warnings := state.Warnings
	state.Mu.Unlock()

	totalWarnings := len(warnings)

	var body strings.Builder
	body.WriteString(styles.BoldAccent.Render(" Scanning..."))
	body.WriteString("\n\n")

	itemsStr := commaFormat(total)
	fmt.Fprintf(&body, " %-13s%s    %s: %s", styles.Dim.Render("Total items:"), styles.Main.Render(itemsStr), styles.Dim.Render("size"), styles.Main.Render(fmt.Sprintf("%s", totalSize)))
	body.WriteString("\n")

	pathStr := currentPath
	maxPath := width - 4
	if len(pathStr) > maxPath {
		pathStr = shortenPath(pathStr, maxPath)
	}
	body.WriteString(" ")

	body.WriteString(styles.Muted.Render(pathStr))
	body.WriteString("\n")

	if totalWarnings > 0 {
		body.WriteString("\n")
		avail := max(height-progressOverhead, 1)
		fmt.Fprintf(&body, " %s", styles.BoldAccent.Render(fmt.Sprintf("WARNING: %d %s - sizes may be inaccurate", totalWarnings, plural("error", totalWarnings))))
		body.WriteString("\n")

		start := 0
		if totalWarnings > avail {
			start = totalWarnings - avail + 1
			fmt.Fprintf(&body, " %s", styles.Muted.Render(fmt.Sprintf("... (%d more) ...", totalWarnings-avail)))
			body.WriteString("\n")
		}
		for i := start; i < totalWarnings; i++ {
			warning := warnings[i]
			maxW := width - 2
			if len(warning) > maxW {
				warning = warning[:maxW]
			}
			fmt.Fprintf(&body, " %s", styles.Dim.Render(warning))
			body.WriteString("\n")
		}
		body.WriteString("\n")
	} else {
		body.WriteString("\n")
	}

	body.WriteString(styles.Muted.Render(" Press ctrl+c to abort"))

	lines := strings.Split(body.String(), "\n")
	for i, line := range lines {
		if len(line) < width {
			lines[i] = line + strings.Repeat(" ", width-len(line))
		}
	}
	return strings.Join(lines, "\n")
}

// plural returns word when count is exactly 1, otherwise word with an "s" appended. Used for
// "1 error" vs "2 errors".
func plural(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

// commaFormat formats an int64 with comma thousands separators. "12345" becomes "12,345".
// Values with three or fewer digits are returned as is.
func commaFormat(value int64) string {
	digits := fmt.Sprintf("%d", value)
	if len(digits) <= commaGroupSize {
		return digits
	}
	var result []byte
	for i, digit := range digits {
		if i > 0 && (len(digits)-i)%commaGroupSize == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(digit))
	}
	return string(result)
}

// shortenPath replaces the middle of a long path with "..." so the start and the end of the
// path remain visible. Returns the path unchanged when it is short enough, or a truncated
// prefix when maxLen is below the minimum the ellipsis needs.
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
