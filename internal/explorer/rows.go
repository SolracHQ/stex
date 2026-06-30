package explorer

import (
	"fmt"
	"math"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	stexmodel "github.com/SolracHQ/stex/internal/model"

	"charm.land/bubbles/v2/table"
)

// Column widths. The name column is widened to fit the terminal so long filenames do not get
// truncated prematurely.
const (
	defaultNameWidth = 40
	sizePctWidth     = 9
	sizeWidth        = 12
)

// buildColumns writes the three table columns, with up or down arrows on the active sort
// column. The name column is widened to fit the terminal so long filenames do not get
// truncated.
func buildColumns(ctx *core.Context) {
	cols := ctx.Table.Columns()
	nameWidth := defaultNameWidth
	if len(cols) > 2 && cols[2].Width > 0 {
		nameWidth = cols[2].Width
	}
	if w := ctx.Width - 2; w > 0 {
		if fit := w - sizePctWidth - sizeWidth - 2; fit > nameWidth {
			nameWidth = fit
		}
	}

	var sizeLabel, nameLabel string
	switch ctx.Config.SortBy {
	case config.SortBySize:
		if ctx.Config.SortOrder == config.Descending {
			sizeLabel = "Size↓"
		} else {
			sizeLabel = "Size↑"
		}
		nameLabel = "Name"
	case config.SortByName:
		sizeLabel = "Size"
		if ctx.Config.SortOrder == config.Descending {
			nameLabel = "Name↓"
		} else {
			nameLabel = "Name↑"
		}
	}

	ctx.Table.SetColumns([]table.Column{
		{Title: " Size%", Width: sizePctWidth},
		{Title: sizeLabel, Width: sizeWidth},
		{Title: nameLabel, Width: nameWidth},
	})
}

// buildRows turns ctx.Items into table rows, with the gradient color and optional emoji
// applied per row.
func buildRows(ctx *core.Context, items []stexmodel.FileSystemItem) {
	rows := make([]table.Row, 0, len(items))
	for _, item := range items {
		rows = append(rows, itemToRow(item, ctx.Config.ShowIcons, ctx.Current))
	}
	ctx.Table.SetRows(rows)
}

// itemToRow converts a single FileSystemItem into a table row. Files and directories get the
// gradient color and the size column, the up link is rendered as a static ".." placeholder.
func itemToRow(item stexmodel.FileSystemItem, showIcons bool, parent *stexmodel.Dir) table.Row {
	switch item.(type) {
	case *stexmodel.File, *stexmodel.Dir:
		var parentSize stexmodel.Size
		if parent != nil {
			parentSize = parent.Size()
		}
		return buildRow(item.Name(), item.Icon(), item.Size(), parentSize, showIcons)
	case *stexmodel.UpLink:
		return table.Row{"", "", "   ..  "}
	}
	return table.Row{}
}

// buildRow formats a single row with the gradient color and the optional emoji prefix on the
// name.
func buildRow(name, emoji string, size, parentSize stexmodel.Size, showIcons bool) table.Row {
	percent := size.PercentOf(parentSize)
	gradientCode := gradientANSI(percent)
	if showIcons {
		name = emoji + " " + name
	}
	return table.Row{
		gradientCode + fmt.Sprintf("%5.2f%%", percent) + "\033[39m",
		gradientCode + " " + size.String() + " \033[39m",
		" " + name + " ",
	}
}

// gradientANSI maps a 0 to 100 percentage to a green to red ANSI foreground color. Values
// outside the range are clamped.
func gradientANSI(percent float64) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	var red, green int
	if percent <= 50 {
		factor := percent / 50.0
		red = int(math.Round(255 * factor))
		green = 255
	} else {
		factor := (percent - 50) / 50.0
		red = 255
		green = int(math.Round(255 * (1 - factor)))
	}
	return fmt.Sprintf("\033[38;2;%d;%d;0m", red, green)
}
