package model

import (
	"sort"

	"github.com/SolracHQ/stex/internal/config"
)

// ComputeItems builds the list of items shown in the directory according to the user's
// grouping, sort, and filter settings. The result always starts with the up link entry when
// the directory has a parent.
func (dir *Dir) ComputeItems(cfg config.Config) []FileSystemItem {
	var items []FileSystemItem

	switch cfg.Grouping {
	case config.FilesFirst:
		items = append(items, sortedSlice(dir.Files(), cfg)...)
		items = append(items, sortedSlice(dir.Dirs(), cfg)...)
	case config.DirsFirst:
		items = append(items, sortedSlice(dir.Dirs(), cfg)...)
		items = append(items, sortedSlice(dir.Files(), cfg)...)
	case config.FilesOnly:
		items = append(items, sortedSlice(dir.Files(), cfg)...)
	case config.DirsOnly:
		items = append(items, sortedSlice(dir.Dirs(), cfg)...)
	case config.Mixed:
		for _, file := range dir.Files() {
			items = append(items, file)
		}
		for _, sub := range dir.Dirs() {
			items = append(items, sub)
		}
		sort.Slice(items, func(i, j int) bool {
			less := lessItem(items[i], items[j], cfg.SortBy)
			if cfg.SortOrder == config.Descending {
				return !less
			}
			return less
		})
	}

	if dir.ParentDir() != nil {
		items = append([]FileSystemItem{NewUpLink(dir.ParentDir())}, items...)
	}

	items = filterItems(items, cfg)

	return items
}

func filterItems(items []FileSystemItem, cfg config.Config) []FileSystemItem {
	if cfg.ShowHidden && cfg.Filter == nil {
		return items
	}
	filtered := make([]FileSystemItem, 0, len(items))
	for _, item := range items {
		if _, ok := item.(*UpLink); ok {
			filtered = append(filtered, item)
			continue
		}
		if !cfg.ShowHidden && len(item.Name()) > 0 && item.Name()[0] == '.' {
			continue
		}
		if cfg.Filter != nil && !cfg.Filter.MatchString(item.Name()) {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}

func lessItem(a, b FileSystemItem, by config.SortBy) bool {
	switch by {
	case config.SortByName:
		nameA, nameB := a.Name(), b.Name()
		if nameA != nameB {
			return nameA < nameB
		}
		return a.UID() < b.UID()
	case config.SortBySize:
		sizeA, sizeB := a.Size(), b.Size()
		if sizeA != sizeB {
			return sizeA < sizeB
		}
		return a.UID() < b.UID()
	}
	return false
}

func sortedSlice[T FileSystemItem](m map[string]T, cfg config.Config) []FileSystemItem {
	out := make([]FileSystemItem, 0, len(m))
	for _, value := range m {
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool {
		r := lessItem(out[i], out[j], cfg.SortBy)
		if cfg.SortOrder == config.Descending {
			return !r
		}
		return r
	})
	return out
}
