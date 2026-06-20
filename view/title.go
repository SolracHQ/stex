package view

import (
	"fmt"

	"github.com/SolracHQ/stex/model"
)

// Title renders the title bar line showing the directory path, file/directory
// counts, and total size. It truncates to width if needed.
func Title(dir *model.Dir, width int) string {
	count := dir.Count()
	text := fmt.Sprintf(" %s - %d files | %d dirs - Total size: %s ",
		dir.FullPath, count.Files, count.Dirs, dir.Size())
	if len(text) > width {
		text = text[:width]
	}
	return text
}
