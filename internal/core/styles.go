package core

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

// Style constants. The colors are chosen to match the original stex terminal look: neutral
// border, cyan help keys, white help descriptions, gray separators.
var (
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	DimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	HelpKeyStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	HelpDescStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	HelpSepStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// TableStyles returns the shared table style. The header is bold, cells are unstyled, the
// selected row uses a dark gray background. This matches the look of the explorer table.
func TableStyles() table.Styles {
	style := table.DefaultStyles()
	style.Header = lipgloss.NewStyle().Bold(true)
	style.Cell = lipgloss.NewStyle()
	style.Selected = lipgloss.NewStyle().Background(lipgloss.Color("239"))
	return style
}
