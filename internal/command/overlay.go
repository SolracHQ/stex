package command

import (
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/styles"
)

const overlayPadX = 8

func (cmd *Command) Overlay(ctx *core.Context) string {
	width := max(1, min(60, ctx.Width-4))
	cmd.input.SetWidth(width - overlayPadX)

	s := cmd.input.Styles()
	s.Focused.Text = styles.Accent
	s.Focused.Prompt = styles.Accent
	cmd.input.SetStyles(s)

	return styles.DialogBorder.Width(width).Render(cmd.input.View())
}
