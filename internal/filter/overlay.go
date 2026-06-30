package filter

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/SolracHQ/stex/internal/core"
)

// Layout constants for the filter overlay. The overlay is a short textinput framed by a single
// border, with horizontal padding so the cursor does not sit on the edge.
const (
	overlayMaxWidth = 60
	overlayHeight   = 5
	overlayPadX     = 8
	liveColor       = "14"
	manualColor     = "11"
)

// overlayBox is the shared style for centred input overlays. The double border in yellow
// matches every other dialog.
var overlayBox = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	BorderForeground(lipgloss.Color("11")).
	Padding(0, 2)

// Overlay renders the filter dialog. The text color is cyan in live mode and yellow in manual
// mode so the user can tell at a glance which mode they are in.
func Overlay(ctx *core.Context, input textinput.Model) string {
	width := max(1, min(overlayMaxWidth, ctx.Width-4))
	input.SetWidth(width - overlayPadX)

	textColor := lipgloss.Color(liveColor)
	if !ctx.Config.LiveFilter {
		textColor = lipgloss.Color(manualColor)
	}
	style := input.Styles()
	style.Focused.Text = lipgloss.NewStyle().Foreground(textColor)
	style.Focused.Prompt = lipgloss.NewStyle().Foreground(textColor)
	input.SetStyles(style)

	return overlayBox.Render(input.View())
}
