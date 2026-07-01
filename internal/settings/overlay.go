package settings

import (
	"fmt"
	"strings"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/styles"
)

func (s *Settings) Overlay(ctx *core.Context) string {
	width := max(1, min(60, ctx.Width-4))
	rows := renderRows(&ctx.Config, s.cursor)

	content := strings.Join([]string{
		styles.BoldAccent.Render("Settings"),
		"",
		rows,
	}, "\n")

	return styles.DialogBorder.Width(width).Render(content)
}

func renderRows(cfg *config.Config, cursor int) string {
	rows := []struct {
		name  string
		value string
	}{
		{"sort", sortLabel(cfg.SortBy)},
		{"order", orderLabel(cfg.SortOrder)},
		{"group", config.GroupingString(cfg.Grouping)},
		{"icons", boolLabel(cfg.ShowIcons, "off", "on")},
		{"hidden", boolLabel(cfg.ShowHidden, "off", "on")},
		{"live filter", boolLabel(cfg.LiveFilter, "off", "on")},
	}
	var b strings.Builder
	for i, r := range rows {
		marker := "  "
		nameStyle := styles.Muted
		if i == cursor {
			marker = styles.BoldAccent.Render("▶ ")
			nameStyle = styles.BoldAccent
		}
		name := nameStyle.Render(padRight(r.name, 12))
		value := styles.Main.Render(r.value)
		b.WriteString(fmt.Sprintf("%s%s  %s", marker, name, value))
		b.WriteString("\n")
	}
	return strings.TrimRight(b.String(), "\n")
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}

func sortLabel(s config.SortBy) string {
	if s == config.SortByName {
		return "name"
	}
	return "size"
}

func orderLabel(o config.SortOrder) string {
	if o == config.Ascending {
		return "asc"
	}
	return "desc"
}

func boolLabel(b bool, off, on string) string {
	if b {
		return on
	}
	return off
}
