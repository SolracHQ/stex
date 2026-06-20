package view

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/model"

	"charm.land/lipgloss/v2"
)

var (
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
)

// DisplayFunc renders a single TreeItem as a display string.
type DisplayFunc func(item model.TreeItem) string

// DefaultDisplay formats a TreeItem as "percent size name".
func DefaultDisplay(item model.TreeItem) string {
	switch item.Kind {
	case model.TKFile:
		percent := item.File.Size.PercentOf(item.File.Parent.Size())
		return fmt.Sprintf("%5s%% %10s %s", percent, item.File.Size, item.File.Name)
	case model.TKDir:
		percent := item.Dir.Size().PercentOf(item.Dir.Parent.Size())
		return fmt.Sprintf("%5s%% %10s %s", percent, item.Dir.Size(), item.Dir.Name)
	case model.TKUpLink:
		return fmt.Sprintf("%17s ..", item.Parent.Size())
	}
	return ""
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
	var buf strings.Builder

	for i, raw := range items[start:end] {
		text := display(raw)
		if width > 0 && len(text) > width {
			text = text[:width]
		}
		if start+i == cursor {
			buf.WriteString(selectedStyle.Render(text))
		} else {
			buf.WriteString(normalStyle.Render(text))
		}
		if i < end-start-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}
