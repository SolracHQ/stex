package model

import (
	"net/http"
	"os"
	"path/filepath"
)

// FileInfo holds metadata about a file or directory for the right pane.
type FileInfo struct {
	Name         string
	Extension    string
	Size         Size
	ModTime      string
	Permissions  string
	IsSymlink    bool
	SymlinkTarget string
	MimeType     string
	IsDir        bool
}

// NewFileInfo gathers metadata for the given path. Returns nil on error.
func NewFileInfo(path string) *FileInfo {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	p := &FileInfo{
		Name:        filepath.Base(path),
		Extension:   filepath.Ext(path),
		Size:        Size(info.Size()),
		ModTime:     info.ModTime().Format("02 Jan 2006 15:04"),
		Permissions: info.Mode().String(),
		IsDir:       info.IsDir(),
	}

	if info.Mode()&os.ModeSymlink != 0 {
		p.IsSymlink = true
		target, err := os.Readlink(path)
		if err == nil {
			p.SymlinkTarget = target
		}
	}

	if !info.IsDir() && info.Mode().IsRegular() {
		p.MimeType = detectMIME(path)
	}

	return p
}

// detectMIME reads the first 512 bytes to detect the content type.
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

// ChildrenInfo holds the top children of a directory for the right pane.
type ChildrenInfo struct {
	Items     []TreeItem
	TotalSize Size
}

// NewChildrenInfo returns up to n largest children of dir.
func NewChildrenInfo(dir *Dir, n int) *ChildrenInfo {
	if dir == nil || n <= 0 {
		return nil
	}
	items := ComputeItems(dir, Config{SortBy: SortBySize, SortOrder: Descending, Grouping: Mixed})
	filtered := make([]TreeItem, 0, len(items))
	for _, item := range items {
		if item.Kind != TKUpLink {
			filtered = append(filtered, item)
		}
	}
	items = filtered
	if len(items) > n {
		items = items[:n]
	}
	var total Size
	for _, item := range items {
		total += item.Size()
	}
	return &ChildrenInfo{Items: items, TotalSize: total}
}


