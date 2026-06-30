package command

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/SolracHQ/stex/internal/core"
)

const (
	overlayMaxWidth = 60
	overlayPadX     = 8
)

func Overlay(ctx *core.Context, input textinput.Model) string {
	width := max(1, min(overlayMaxWidth, ctx.Width-4))
	input.SetWidth(width - overlayPadX)

	style := input.Styles()
	style.Focused.Text = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	style.Focused.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	input.SetStyles(style)

	return lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("11")).
		Padding(0, 2).
		Width(width).
		Render(input.View())
}
