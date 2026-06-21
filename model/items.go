package model

import "sort"

// TreeKind classifies items in a directory listing.
type TreeKind int

const (
	TKFile   TreeKind = iota
	TKDir
	TKUpLink
)

// TreeItem is a single row in the file listing.
type TreeItem struct {
	Kind   TreeKind
	File   *File
	Dir    *Dir
	Parent *Dir // set for TKUpLink
}

// Name returns the display name of the item.
func (ti TreeItem) Name() string {
	switch ti.Kind {
	case TKFile:
		return ti.File.Name
	case TKDir:
		return ti.Dir.Name
	}
	return ""
}

// Size returns the size of the item.
func (ti TreeItem) Size() Size {
	switch ti.Kind {
	case TKFile:
		return ti.File.Size
	case TKDir:
		return ti.Dir.Size()
	}
	return 0
}

// Icon returns the emoji for the item kind.
func (ti TreeItem) Icon() string {
	switch ti.Kind {
	case TKFile:
		return "📄"
	case TKDir:
		return "📁"
	}
	return ""
}

// FullPath returns the absolute path of the item.
func (ti TreeItem) FullPath() string {
	switch ti.Kind {
	case TKFile:
		return ti.File.FullPath
	case TKDir:
		return ti.Dir.FullPath
	}
	return ""
}

// ParentDir returns the parent directory of the item.
func (ti TreeItem) ParentDir() *Dir {
	switch ti.Kind {
	case TKFile:
		return ti.File.Parent
	case TKDir:
		return ti.Dir.Parent
	case TKUpLink:
		return ti.Parent
	}
	return nil
}

// ComputeItems returns sorted, grouped items for dir. Prepends an up-link
// entry when dir has a parent.
func ComputeItems(dir *Dir, cfg Config) []TreeItem {
	var items []TreeItem

	switch cfg.Grouping {
	case FilesFirst:
		items = appendFileItems(items, sortedFileSlice(dir.Files, cfg))
		items = appendDirItems(items, sortedDirSlice(dir.Dirs, cfg))
	case DirsFirst:
		items = appendDirItems(items, sortedDirSlice(dir.Dirs, cfg))
		items = appendFileItems(items, sortedFileSlice(dir.Files, cfg))
	case FilesOnly:
		items = appendFileItems(items, sortedFileSlice(dir.Files, cfg))
	case DirsOnly:
		items = appendDirItems(items, sortedDirSlice(dir.Dirs, cfg))
	case Mixed:
		items = appendFileItems(items, dir.Files)
		items = appendDirItems(items, dir.Dirs)
		sort.Slice(items, func(i, j int) bool {
			less := lessItem(items[i], items[j], cfg.SortBy)
			if cfg.SortOrder == Descending {
				return !less
			}
			return less
		})
	}

	if dir.Parent != nil {
		items = append([]TreeItem{{Kind: TKUpLink, Parent: dir.Parent}}, items...)
	}

	return items
}

func lessFile(a File, b File, by SortBy) bool {
	switch by {
	case SortByName:
		return a.Name < b.Name
	case SortBySize:
		return a.Size < b.Size
	}
	return false
}

func lessDir(a *Dir, b *Dir, by SortBy) bool {
	switch by {
	case SortByName:
		return a.Name < b.Name
	case SortBySize:
		return a.Size() < b.Size()
	}
	return false
}

func lessItem(a TreeItem, b TreeItem, by SortBy) bool {
	switch by {
	case SortByName:
		return a.Name() < b.Name()
	case SortBySize:
		return a.Size() < b.Size()
	}
	return false
}

func sortedFileSlice(files []File, cfg Config) []File {
	out := make([]File, len(files))
	copy(out, files)
	sort.Slice(out, func(i, j int) bool {
		less := lessFile(out[i], out[j], cfg.SortBy)
		if cfg.SortOrder == Descending {
			return !less
		}
		return less
	})
	return out
}

func sortedDirSlice(dirs []*Dir, cfg Config) []*Dir {
	out := make([]*Dir, len(dirs))
	copy(out, dirs)
	sort.Slice(out, func(i, j int) bool {
		less := lessDir(out[i], out[j], cfg.SortBy)
		if cfg.SortOrder == Descending {
			return !less
		}
		return less
	})
	return out
}

func appendFileItems(items []TreeItem, files []File) []TreeItem {
	for i := range files {
		items = append(items, TreeItem{Kind: TKFile, File: &files[i]})
	}
	return items
}

func appendDirItems(items []TreeItem, dirs []*Dir) []TreeItem {
	for _, d := range dirs {
		items = append(items, TreeItem{Kind: TKDir, Dir: d})
	}
	return items
}
