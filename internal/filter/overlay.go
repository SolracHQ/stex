package filter

import (
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/styles"
)

// Overlay renders the filter dialog centered on the base view. The text color is cyan in live
// mode and yellow in manual mode so the user can tell at a glance which mode they are in.
func (flt *Filter) Overlay(ctx *core.Context) string {
	width := max(1, min(60, ctx.Width-4))
	flt.input.SetWidth(width - 8)

	textStyle := styles.Accent
	if ctx.Config.LiveFilter {
		textStyle = styles.Active
	}
	s := flt.input.Styles()
	s.Focused.Text = textStyle
	s.Focused.Prompt = textStyle
	flt.input.SetStyles(s)

	return styles.DialogBorder.Width(width).Render(flt.input.View())
}
