package view

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/SolracHQ/stex/model"

	"charm.land/lipgloss/v2"
)

var (
	normalStyle   = lipgloss.NewStyle()
	selectedStyle = lipgloss.NewStyle().Background(lipgloss.Color("239"))
)

// gradientHex returns a hex color string for a percentage 0-100, mapping
// green (0%) through yellow (50%) to red (100%).
func gradientHex(percentage float64) string {
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}
	var red, green int
	if percentage <= 50 {
		factor := percentage / 50.0
		red = int(math.Round(255 * factor))
		green = 255
	} else {
		factor := (percentage - 50) / 50.0
		red = 255
		green = int(math.Round(255 * (1 - factor)))
	}
	return fmt.Sprintf("#%02x%02x00", red, green)
}

func gradientStyle(percentage float64) lipgloss.Style {
	if percentage < 0 {
		return normalStyle
	}
	return lipgloss.NewStyle().Foreground(lipgloss.Color(gradientHex(percentage)))
}

// DisplayFunc renders a single TreeItem as a display string, returning the
// percentage value for gradient colouring (negative = no colouring).
type DisplayFunc func(item model.TreeItem) (percentage float64, text string)

// DefaultDisplay formats a TreeItem as "percent size name".
func DefaultDisplay(item model.TreeItem) (float64, string) {
	switch item.Kind {
	case model.TKFile:
		percentage := item.File.Size.PercentOf(item.File.Parent.Size())
		parsedFloat, _ := strconv.ParseFloat(percentage, 64)
		return parsedFloat, fmt.Sprintf("%5s%% %10s %s", percentage, item.File.Size, item.File.Name)
	case model.TKDir:
		percentage := item.Dir.Size().PercentOf(item.Dir.Parent.Size())
		parsedFloat, _ := strconv.ParseFloat(percentage, 64)
		return parsedFloat, fmt.Sprintf("%5s%% %10s %s", percentage, item.Dir.Size(), item.Dir.Name)
	case model.TKUpLink:
		return -1, fmt.Sprintf("%18s..", "")
	}
	return -1, ""
}

// visibleRange returns the item slice that fits in available rows, keeping
// the cursor centred.
func visibleRange(index, size, available int) (start, end int) {
	if size <= available {
		return 0, size
	}
	half := available / 2
	switch {
	case index <= half:
		return 0, available
	case index >= size-half-1:
		return size - available, size
	default:
		return index - half, index + half + 1
	}
}

// List renders a scrollable, selectable list of TreeItems.
func List(items []model.TreeItem, cursor, height, width int, display DisplayFunc) string {
	if len(items) == 0 || height <= 0 {
		return ""
	}

	start, end := visibleRange(cursor, len(items), height)
	var buffer strings.Builder

	for i, raw := range items[start:end] {
		percentage, text := display(raw)
		if width > 0 && len(text) > width {
			text = text[:width]
		}

		if start+i == cursor {
			if width > 0 && len(text) < width {
				text += strings.Repeat(" ", width-len(text))
			}
			if percentage >= 0 {
				text = gradientStyle(percentage).Render(text)
			}
			text = selectedStyle.Render(text)
			buffer.WriteString(text)
		} else {
			if percentage >= 0 {
				buffer.WriteString(gradientStyle(percentage).Render(text))
			} else {
				buffer.WriteString(text)
			}
		}
		if i < end-start-1 {
			buffer.WriteString("\n")
		}
	}

	return buffer.String()
}
