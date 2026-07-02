// Package styles centralizes the visual constants and reusable style builders for every
// dialog, input, and display element in the application.
package styles

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/table"
	"charm.land/lipgloss/v2"
)

// Color palette.
const (
	AccentColor = "11"  // yellow, used for borders, highlights, cursor, titles
	ActiveColor = "14"  // cyan, used for live indicators and help keys
	MainColor   = "7"   // white, used for primary text
	DimColor    = "240" // gray, used for borders, separators, secondary text
	MutedColor  = "8"   // darker gray, used for less important text
	SelectBg    = "239" // dark gray background for selected table rows
)

// DialogBorder returns a double border style with the accent color and inner horizontal
// padding. The caller should set Width before rendering.
var DialogBorder = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color(AccentColor)).
	Padding(0, 2)

// DialogBorderWide is like DialogBorder but with wider horizontal padding.
var DialogBorderWide = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color(AccentColor)).
	Padding(0, 4)

// BorderNorm is a normal (single line) border with dim foreground.
var BorderNorm = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color(DimColor))

// Dim is a dim foreground style for separators and less prominent text.
var Dim = lipgloss.NewStyle().Foreground(lipgloss.Color(DimColor))

// Accent renders text in accent color (yellow).
var Accent = lipgloss.NewStyle().Foreground(lipgloss.Color(AccentColor))

// Active renders text in active color (cyan).
var Active = lipgloss.NewStyle().Foreground(lipgloss.Color(ActiveColor))

// BoldAccent renders text in accent color and bold.
var BoldAccent = lipgloss.NewStyle().Foreground(lipgloss.Color(AccentColor)).Bold(true)

// Main renders text in main (white) color.
var Main = lipgloss.NewStyle().Foreground(lipgloss.Color(MainColor))

// BoldMain renders text in main color and bold.
var BoldMain = lipgloss.NewStyle().Foreground(lipgloss.Color(MainColor)).Bold(true)

// Muted renders text in muted (darker gray) color.
var Muted = lipgloss.NewStyle().Foreground(lipgloss.Color(MutedColor))

// HelpKey is the style for the key part of the help footer.
var HelpKey = lipgloss.NewStyle().Foreground(lipgloss.Color(ActiveColor))

// HelpDesc is the style for the description part of the help footer.
var HelpDesc = lipgloss.NewStyle().Foreground(lipgloss.Color(MainColor))

// HelpSep is the style for separators in the help footer.
var HelpSep = lipgloss.NewStyle().Foreground(lipgloss.Color(DimColor))

// HelpDefaults returns a help.Model with styles configured to match the application palette.
func HelpDefaults() help.Model {
	m := help.New()
	m.Styles.FullKey = HelpKey
	m.Styles.FullDesc = HelpDesc
	m.Styles.FullSeparator = HelpSep
	m.Styles.ShortKey = HelpKey
	m.Styles.ShortDesc = HelpDesc
	m.Styles.ShortSeparator = HelpSep
	return m
}

// TableDefault returns the standard table styles.
func TableDefault() table.Styles {
	s := table.DefaultStyles()
	s.Header = lipgloss.NewStyle().Bold(true)
	s.Cell = lipgloss.NewStyle()
	s.Selected = lipgloss.NewStyle().Background(lipgloss.Color(SelectBg))
	return s
}
