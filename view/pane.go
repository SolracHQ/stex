package view

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// Split joins left and right column strings with sep between them.
// The shorter side is padded with empty lines. Each line is truncated or
// padded to its column width. ANSI sequences are handled correctly.
func Split(left, right string, leftWidth, rightWidth int, sep string) string {
	leftLines := strings.Split(left, "\n")
	rightLines := strings.Split(right, "\n")

	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	var buf strings.Builder
	for i := 0; i < maxLines; i++ {
		if i > 0 {
			buf.WriteByte('\n')
		}

		var leftLine, rightLine string
		if i < len(leftLines) {
			leftLine = leftLines[i]
		}
		if i < len(rightLines) {
			rightLine = rightLines[i]
		}

		leftLine = ansi.Truncate(leftLine, leftWidth, "")
		if ansi.StringWidth(leftLine) < leftWidth {
			leftLine += strings.Repeat(" ", leftWidth-ansi.StringWidth(leftLine))
		}

		rightLine = ansi.Truncate(rightLine, rightWidth, "")
		if ansi.StringWidth(rightLine) < rightWidth {
			rightLine += strings.Repeat(" ", rightWidth-ansi.StringWidth(rightLine))
		}

		buf.WriteString(leftLine)
		buf.WriteString(sep)
		buf.WriteString(rightLine)
	}

	return buf.String()
}
