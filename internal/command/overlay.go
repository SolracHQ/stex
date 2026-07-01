package command

import (
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/styles"
)

const overlayPadX = 8

func (c *Command) Overlay(ctx *core.Context) string {
	width := max(1, min(60, ctx.Width-4))
	c.input.SetWidth(width - overlayPadX)

	s := c.input.Styles()
	s.Focused.Text = styles.Accent
	s.Focused.Prompt = styles.Accent
	c.input.SetStyles(s)

	return styles.DialogBorder.Width(width).Render(c.input.View())
}
