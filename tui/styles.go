package tui

import "charm.land/lipgloss/v2"

var (
	// scanBorderStyle is the outer border used around the scan progress dialog.
	scanBorderStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

	// borderStyle is the outer border used around the main file listing.
	borderStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))

	// helpKeyStyle, helpDescStyle, and helpSepStyle style the key-binding
	// display produced by the bubbles/help model.
	helpKeyStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	helpDescStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	helpSepStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	// dialogBoxStyle wraps the help overlay in a rounded, cyan-bordered box.
	dialogBoxStyle    = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("14")).Padding(1, 2)
	dialogTitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
	dialogFooterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)
