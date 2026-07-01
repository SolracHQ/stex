package filter

import (
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/styles"
)

// Overlay renders the filter dialog centered on the base view. The text color is cyan in live
// mode and yellow in manual mode so the user can tell at a glance which mode they are in.
func (f *Filter) Overlay(ctx *core.Context) string {
	width := max(1, min(60, ctx.Width-4))
	f.input.SetWidth(width - 8)

	textStyle := styles.Accent
	if ctx.Config.LiveFilter {
		textStyle = styles.Active
	}
	s := f.input.Styles()
	s.Focused.Text = textStyle
	s.Focused.Prompt = textStyle
	f.input.SetStyles(s)

	return styles.DialogBorder.Width(width).Render(f.input.View())
}
