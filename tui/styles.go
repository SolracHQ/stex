package tui

import (
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

var (
	scanBorderStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	borderStyle     = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	helpKeyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	helpDescStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	helpSepStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	dialogBoxStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("14")).Padding(1, 2)
	dialogTitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
	dialogFooterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func tableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Header = lipgloss.NewStyle().Bold(true)
	s.Cell = lipgloss.NewStyle()
	s.Selected = lipgloss.NewStyle().Background(lipgloss.Color("239"))
	return s
}
