package view

import (
	"fmt"

	"github.com/SolracHQ/stex/model"
)

// ColumnHeaders renders column labels with a sort-direction arrow.
func ColumnHeaders(cfg model.Config, width int) string {
	var (
		sizeHead = fmt.Sprintf("%10s", "Size")
		nameHead = "Name"
	)

	switch cfg.SortBy {
	case model.SortBySize:
		if cfg.SortOrder == model.Descending {
			sizeHead = fmt.Sprintf("%9s↓", "Size")
		} else {
			sizeHead = fmt.Sprintf("%9s↑", "Size")
		}
	case model.SortByName:
		if cfg.SortOrder == model.Descending {
			nameHead = "Name↓"
		} else {
			nameHead = "Name↑"
		}
	}

	text := fmt.Sprintf("%6s %10s %s", "Size%", sizeHead, nameHead)
	if len(text) > width {
		text = text[:width]
	}
	return text
}
