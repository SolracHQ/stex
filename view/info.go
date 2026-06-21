package view

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/model"
)

// RenderFileInfo returns a formatted metadata panel for a file.
func RenderFileInfo(info *model.FileInfo, width, height int) string {
	if info == nil {
		return centered("No info available", width, height)
	}

	var lines []string
	lines = append(lines, bold(" File Info "))
	lines = append(lines, "")

	lines = append(lines, fmt.Sprintf(" %-14s%s", "Name:", info.Name))
	if info.Extension != "" {
		lines = append(lines, fmt.Sprintf(" %-14s%s", "Extension:", info.Extension))
	}
	if info.MimeType != "" {
		lines = append(lines, fmt.Sprintf(" %-14s%s", "Type:", info.MimeType))
	}
	lines = append(lines, fmt.Sprintf(" %-14s%s", "Size:", info.Size))
	lines = append(lines, fmt.Sprintf(" %-14s%s", "Modified:", info.ModTime))
	lines = append(lines, fmt.Sprintf(" %-14s%s", "Permissions:", info.Permissions))
	if info.IsSymlink {
		lines = append(lines, fmt.Sprintf(" %-14s%s", "Symlink:", info.SymlinkTarget))
	}

	return padLines(lines, width, height)
}

// RenderDirInfo returns a formatted metadata panel for a directory,
// including the top children sorted by size.
func RenderDirInfo(dir *model.Dir, info *model.FileInfo, children *model.ChildrenInfo, width, height int) string {
	if info == nil && dir == nil {
		return centered("No info available", width, height)
	}

	var lines []string
	lines = append(lines, bold(" Directory Info "))
	lines = append(lines, "")

	if info != nil {
		lines = append(lines, fmt.Sprintf(" %-14s%s", "Name:", info.Name))
		dirSize := info.Size
		if dir != nil {
			dirSize = dir.Size()
		}
		lines = append(lines, fmt.Sprintf(" %-14s%s", "Size:", dirSize))
		lines = append(lines, fmt.Sprintf(" %-14s%s", "Modified:", info.ModTime))
		lines = append(lines, fmt.Sprintf(" %-14s%s", "Permissions:", info.Permissions))
	}

	if children != nil && len(children.Items) > 0 {
		lines = append(lines, "")
		lines = append(lines, bold(" Largest Children "))
		lines = append(lines, "")
		for _, item := range children.Items {
			pct := itemSizePercent(item, children.TotalSize)
			lines = append(lines, fmt.Sprintf(" %5.2f%% %10s  %s", pct, itemDisplaySize(item), item.Name()))
		}
	}

	return padLines(lines, width, height)
}

func centered(text string, width, height int) string {
	lines := strings.Split(text, "\n")
	contentHeight := len(lines)
	topPad := (height - contentHeight) / 2
	if topPad < 0 {
		topPad = 0
	}
	var buf strings.Builder
	for i := 0; i < topPad; i++ {
		buf.WriteString("\n")
	}
	for _, line := range lines {
		padding := (width - len(line)) / 2
		if padding < 0 {
			padding = 0
		}
		buf.WriteString(strings.Repeat(" ", padding))
		buf.WriteString(line)
		buf.WriteString("\n")
	}
	return strings.TrimRight(buf.String(), "\n")
}

func padLines(lines []string, width, height int) string {
	var buf strings.Builder
	for i, line := range lines {
		if i > 0 {
			buf.WriteString("\n")
		}
		if len(line) > width {
			line = line[:width]
		}
		buf.WriteString(line)
	}
	used := len(lines)
	for i := used; i < height; i++ {
		buf.WriteString("\n")
	}
	return buf.String()
}

func bold(text string) string {
	return "\033[1m" + text + "\033[22m"
}

func itemSizePercent(item model.TreeItem, total model.Size) float64 {
	s := item.Size()
	if total == 0 {
		return 0
	}
	return float64(s) / float64(total) * 100
}

func itemDisplaySize(item model.TreeItem) string {
	return item.Size().String()
}
