package core

import (
	"strings"

	"charm.land/bubbles/v2/help"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// SplitViewThreshold is the minimum inner width at which the layout switches from a single
// column to the table plus info side by side.
const SplitViewThreshold = 80

// Blank returns a height by width canvas of spaces. Used as the base behind centred overlay
// dialogs during the initial scan and as the body of any mode that wants to draw an empty
// backdrop. Returns "" when either dimension is non positive.
func Blank(width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	emptyLine := strings.Repeat(" ", width)
	lines := make([]string, height)
	for index := range lines {
		lines[index] = emptyLine
	}
	return strings.Join(lines, "\n")
}

// RenderBase renders the base panels, title, table, info, and footer, using the given help key
// map. The mode that owns the active state passes its key list here so the footer reflects
// what it can do. Returns "" when the context has no dimensions or no current directory.
func RenderBase(ctx *Context, keys help.KeyMap) string {
	if ctx.Width == 0 || ctx.Height == 0 || ctx.Current == nil {
		return ""
	}
	innerWidth := ctx.Width - 2
	innerHeight := ctx.Height - 2

	if innerWidth >= SplitViewThreshold {
		return splitView(ctx, innerWidth, innerHeight, keys)
	}
	return narrowView(ctx, innerWidth, innerHeight, keys)
}

// splitView renders the dual pane layout. The left half is the file table, the right half is
// the info pane, separated by a vertical bar.
func splitView(ctx *Context, innerWidth, innerHeight int, keys help.KeyMap) string {
	leftWidth := innerWidth/2 - 1
	rightWidth := innerWidth - leftWidth - 1

	titleLine := Title(ctx.Current, innerWidth, ctx.Config.ShowIcons, TitleGroup(ctx))
	sepLine := DimStyle.Render(strings.Repeat("─", innerWidth))

	footerStr, footerHeight := renderFooter(ctx, innerWidth, keys)
	contentHeight := max(innerHeight-2-footerHeight, 0)

	leftContent := renderPadded(ctx.Table.View(), contentHeight)

	var rightContent string
	if ctx.Info.Content != "" {
		rightContent = renderPadded(ctx.Info.Content, contentHeight)
	}

	combined := splitPane(leftContent, rightContent, leftWidth, rightWidth, DimStyle.Render("│"))

	return titleLine + "\n" + sepLine + "\n" + combined + "\n" + strings.Join(footerStr, "\n")
}

// narrowView renders the single column layout used on terminals narrower than the split
// threshold. The info pane is hidden.
func narrowView(ctx *Context, innerWidth, innerHeight int, keys help.KeyMap) string {
	footerStr, footerHeight := renderFooter(ctx, innerWidth, keys)
	contentHeight := max(innerHeight-2-footerHeight, 0)

	lines := make([]string, innerHeight)
	lines[0] = Title(ctx.Current, innerWidth, ctx.Config.ShowIcons, TitleGroup(ctx))
	lines[1] = DimStyle.Render(strings.Repeat("─", innerWidth))

	fillTable(lines, contentHeight, ctx.Table.View())

	footerStart := 2 + contentHeight
	for i := 0; i < footerHeight && footerStart+i < innerHeight; i++ {
		lines[footerStart+i] = footerStr[i]
	}
	for i := footerStart + footerHeight; i < innerHeight; i++ {
		lines[i] = ""
	}

	return strings.Join(lines, "\n")
}

// renderFooter returns the footer content and its line count. The footer is the centred help
// line, or empty when no keys are supplied. Returning a nil slice when keys is nil keeps
// callers from having to filter empty footers downstream.
func renderFooter(ctx *Context, width int, keys help.KeyMap) ([]string, int) {
	if keys == nil {
		return nil, 0
	}
	content := ctx.Help.View(keys)
	if content == "" {
		return nil, 0
	}
	content = lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(content)
	lines := strings.Split(content, "\n")
	return lines, len(lines)
}

// fillTable copies the rendered table lines into the content area of the screen buffer,
// starting at line index 2 (after the title and the separator). Empty lines are written past
// the end of the table to keep the content area a fixed size.
func fillTable(lines []string, contentHeight int, tableContent string) {
	tableLines := strings.Split(tableContent, "\n")
	for i := 0; i < contentHeight && i < len(tableLines); i++ {
		lines[2+i] = tableLines[i]
	}
	for i := len(tableLines); i < contentHeight; i++ {
		lines[2+i] = ""
	}
}

// renderPadded pads content to exactly height lines. Lines beyond the source are filled with
// empty strings.
func renderPadded(content string, height int) string {
	lines := strings.Split(content, "\n")
	out := make([]string, height)
	for i := 0; i < height && i < len(lines); i++ {
		out[i] = lines[i]
	}
	for i := len(lines); i < height; i++ {
		out[i] = ""
	}
	return strings.Join(out, "\n")
}

// WrapView applies the standard border and screen settings to the given body and returns the
// resulting tea.View. The alternate screen is enabled and cell motion mouse is turned on, so
// the app captures the mouse for click and wheel handling.
func WrapView(content string) tea.View {
	view := tea.NewView(BorderStyle.Render(content))
	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion
	return view
}

// LoadingView returns a tea.View configured for the loading screen. The body is rendered as
// is, the alternate screen is enabled, and all motion mouse is used so the progress dialog
// updates cleanly while the user can still move the mouse.
func LoadingView(content string) tea.View {
	view := tea.NewView(content)
	view.AltScreen = true
	view.MouseMode = tea.MouseModeAllMotion
	return view
}

// splitPane joins left and right column strings with sep between them, line by line. The
// shorter side is padded with empty lines. Each line is truncated or padded to its column
// width. ANSI escape sequences are handled with the ansi package so colored content is not
// corrupted by the truncation.
func splitPane(left, right string, leftWidth, rightWidth int, sep string) string {
	leftLines := strings.Split(left, "\n")
	rightLines := strings.Split(right, "\n")

	maxLines := max(len(rightLines), len(leftLines))

	var buf strings.Builder
	for index := range maxLines {
		if index > 0 {
			buf.WriteByte('\n')
		}

		var leftLine, rightLine string
		if index < len(leftLines) {
			leftLine = leftLines[index]
		}
		if index < len(rightLines) {
			rightLine = rightLines[index]
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
