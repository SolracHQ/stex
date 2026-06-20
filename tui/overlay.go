package tui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// overlayCenter places foreground over background, centred.
func overlayCenter(background, foreground string) string {
	if foreground == "" || background == "" {
		return background
	}

	fgWidth, fgHeight := lipgloss.Size(foreground)
	bgWidth, bgHeight := lipgloss.Size(background)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		return foreground
	}

	x := (bgWidth - fgWidth) / 2
	y := (bgHeight - fgHeight) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	fgLines := strings.Split(foreground, "\n")
	bgLines := strings.Split(background, "\n")

	var buf strings.Builder
	for i, bgLine := range bgLines {
		if i > 0 {
			buf.WriteByte('\n')
		}
		if i < y || i >= y+fgHeight {
			buf.WriteString(bgLine)
			continue
		}

		pos := 0
		if x > 0 {
			left := ansi.Truncate(bgLine, x, "")
			pos = ansi.StringWidth(left)
			buf.WriteString(left)
			if pos < x {
				buf.WriteString(strings.Repeat(" ", x-pos))
				pos = x
			}
		}

		fgLine := fgLines[i-y]
		buf.WriteString(fgLine)
		pos += ansi.StringWidth(fgLine)

		right := ansi.TruncateLeft(bgLine, pos, "")
		bgW := ansi.StringWidth(bgLine)
		rightW := ansi.StringWidth(right)
		if rightW <= bgW-pos {
			buf.WriteString(strings.Repeat(" ", bgW-rightW-pos))
		}
		buf.WriteString(right)
	}

	return buf.String()
}
