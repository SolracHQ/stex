package core

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/SolracHQ/stex/internal/model"
)

// FileInfo holds the metadata that the right pane shows for a single file. The fields mirror
// the read only stat information plus the MIME type, which is detected from the first 512 bytes
// for regular files.
type FileInfo struct {
	Name          string
	Extension     string
	Size          model.Size
	ModTime       string
	Permissions   string
	IsSymlink     bool
	SymlinkTarget string
	MimeType      string
	IsDir         bool
}

// NewFileInfo reads the metadata of the file at path. Returns nil when the file does not exist
// or is not readable. The returned value is always safe to read, even for the nil case.
func NewFileInfo(path string) *FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	out := &FileInfo{
		Name:        filepath.Base(path),
		Extension:   filepath.Ext(path),
		Size:        model.Size(info.Size()),
		ModTime:     info.ModTime().Format("02 Jan 2006 15:04"),
		Permissions: info.Mode().String(),
		IsDir:       info.IsDir(),
	}

	if info.Mode()&os.ModeSymlink != 0 {
		out.IsSymlink = true
		if target, err := os.Readlink(path); err == nil {
			out.SymlinkTarget = target
		}
	}

	if !info.IsDir() && info.Mode().IsRegular() {
		out.MimeType = detectMIME(path)
	}

	return out
}

// detectMIME reads the first 512 bytes of path and returns the MIME type detected by net/http.
// Returns an empty string if the file cannot be opened or read.
func detectMIME(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	return http.DetectContentType(buf[:n])
}

// ChildrenInfo holds the largest direct children of a directory for the "Largest Children"
// segment of the right pane.
type ChildrenInfo struct {
	Items     []model.FileSystemItem
	TotalSize model.Size
}

// RenderFileInfo renders the right pane for a file. width and height are the cell budget for
// the pane, the result is padded out to fit. Returns a centred "No info available" when info
// is nil.
func RenderFileInfo(info *FileInfo, width, height int) string {
	if info == nil {
		return centered("No info available", width, height)
	}

	var lines []string
	lines = append(lines, bold(" File Info "))
	lines = append(lines, "")

	lines = append(lines, fmt.Sprintf(" %-14s%s", "Name:", info.Name))
	if info.Extension != "" {
		lines = append(lines, fmt.Sprintf(" %-14s%s", "Extension:", info.Extension))
	}
	if info.MimeType != "" {
		lines = append(lines, fmt.Sprintf(" %-14s%s", "Type:", info.MimeType))
	}
	lines = append(lines, fmt.Sprintf(" %-14s%s", "Size:", info.Size))
	lines = append(lines, fmt.Sprintf(" %-14s%s", "Modified:", info.ModTime))
	lines = append(lines, fmt.Sprintf(" %-14s%s", "Permissions:", info.Permissions))
	if info.IsSymlink {
		lines = append(lines, fmt.Sprintf(" %-14s%s", "Symlink:", info.SymlinkTarget))
	}

	return padLines(lines, width, height)
}

// RenderDirInfo renders the right pane for a directory. The total size is taken from dirSize,
// which the caller passes in separately so the size of the directory can come from the model
// (which has it cached) instead of being re stat'd. The children segment shows the top n items
// from children, sorted by size.
func RenderDirInfo(info *FileInfo, dirSize model.Size, children *ChildrenInfo, width, height int) string {
	if info == nil {
		return centered("No info available", width, height)
	}

	var lines []string
	lines = append(lines, bold(" Directory Info "))
	lines = append(lines, "")

	lines = append(lines, fmt.Sprintf(" %-14s%s", "Name:", info.Name))
	lines = append(lines, fmt.Sprintf(" %-14s%s", "Size:", dirSize))
	lines = append(lines, fmt.Sprintf(" %-14s%s", "Modified:", info.ModTime))
	lines = append(lines, fmt.Sprintf(" %-14s%s", "Permissions:", info.Permissions))

	if children != nil && len(children.Items) > 0 {
		lines = append(lines, "")
		lines = append(lines, bold(" Largest Children "))
		lines = append(lines, "")
		for _, item := range children.Items {
			pct := item.Size().PercentOf(children.TotalSize)
			lines = append(lines, fmt.Sprintf(" %5.2f%% %10s  %s", pct, item.Size().String(), item.Name()))
		}
	}

	return padLines(lines, width, height)
}

// centered writes text centred inside a width by height box. Lines are split on "\n" and each
// line is padded with spaces to the centre. Extra vertical space is added above the text.
func centered(text string, width, height int) string {
	lines := strings.Split(text, "\n")
	contentHeight := len(lines)
	topPad := max((height-contentHeight)/2, 0)
	var buf strings.Builder
	for range topPad {
		buf.WriteString("\n")
	}
	for _, line := range lines {
		padding := max((width-len(line))/2, 0)
		buf.WriteString(strings.Repeat(" ", padding))
		buf.WriteString(line)
		buf.WriteString("\n")
	}
	return strings.TrimRight(buf.String(), "\n")
}

// padLines joins the given lines with newlines, truncating any line that exceeds width, then
// appends empty lines until the result has exactly height lines.
func padLines(lines []string, width, height int) string {
	var buf strings.Builder
	for index, line := range lines {
		if index > 0 {
			buf.WriteString("\n")
		}
		if len(line) > width {
			line = line[:width]
		}
		buf.WriteString(line)
	}
	used := len(lines)
	for i := used; i < height; i++ {
		buf.WriteString("\n")
	}
	return buf.String()
}

// bold returns text wrapped in the ANSI escape for bold on and off. Used for the section
// headers in the info pane.
func bold(text string) string {
	return "\033[1m" + text + "\033[22m"
}
