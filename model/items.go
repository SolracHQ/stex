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

// ComputeItems returns sorted, grouped items for dir. Prepends an up-link
// entry when dir has a parent.
func ComputeItems(dir *Dir, cfg Config) []TreeItem {
	numFiles := len(dir.Files)
	numDirs := len(dir.Dirs)

	sortedFiles := make([]File, numFiles)
	copy(sortedFiles, dir.Files)
	sortedDirs := make([]*Dir, numDirs)
	copy(sortedDirs, dir.Dirs)

	sort.Slice(sortedFiles, func(i, j int) bool {
		less := lessFile(sortedFiles[i], sortedFiles[j], cfg.SortBy)
		if cfg.SortOrder == Descending {
			return !less
		}
		return less
	})

	sort.Slice(sortedDirs, func(i, j int) bool {
		less := lessDir(sortedDirs[i], sortedDirs[j], cfg.SortBy)
		if cfg.SortOrder == Descending {
			return !less
		}
		return less
	})

	var items []TreeItem

	switch cfg.Grouping {
	case FilesFirst:
		for i := 0; i < numFiles; i++ {
			items = append(items, TreeItem{Kind: TKFile, File: &sortedFiles[i]})
		}
		for i := 0; i < numDirs; i++ {
			items = append(items, TreeItem{Kind: TKDir, Dir: sortedDirs[i]})
		}
	case DirsFirst:
		for i := 0; i < numDirs; i++ {
			items = append(items, TreeItem{Kind: TKDir, Dir: sortedDirs[i]})
		}
		for i := 0; i < numFiles; i++ {
			items = append(items, TreeItem{Kind: TKFile, File: &sortedFiles[i]})
		}
	case FilesOnly:
		for i := 0; i < numFiles; i++ {
			items = append(items, TreeItem{Kind: TKFile, File: &sortedFiles[i]})
		}
	case DirsOnly:
		for i := 0; i < numDirs; i++ {
			items = append(items, TreeItem{Kind: TKDir, Dir: sortedDirs[i]})
		}
	case Mixed:
		for i := 0; i < numFiles; i++ {
			items = append(items, TreeItem{Kind: TKFile, File: &sortedFiles[i]})
		}
		for i := 0; i < numDirs; i++ {
			items = append(items, TreeItem{Kind: TKDir, Dir: sortedDirs[i]})
		}
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
		return itemName(a) < itemName(b)
	case SortBySize:
		return itemSize(a) < itemSize(b)
	}
	return false
}

func itemName(item TreeItem) string {
	switch item.Kind {
	case TKFile:
		return item.File.Name
	case TKDir:
		return item.Dir.Name
	}
	return ""
}

func itemSize(item TreeItem) Size {
	switch item.Kind {
	case TKFile:
		return item.File.Size
	case TKDir:
		return item.Dir.Size()
	}
	return 0
}
