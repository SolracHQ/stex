package explorer

import (
	"sort"

	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	stexmodel "github.com/SolracHQ/stex/internal/model"
)

// ComputeItems builds the list of items shown in the current directory according to the user's
// grouping, sort, and filter settings. The result always starts with the up link entry when
// the directory has a parent.
func ComputeItems(dir *stexmodel.Dir, cfg config.Config) []stexmodel.FileSystemItem {
	var items []stexmodel.FileSystemItem

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
		items = append([]stexmodel.FileSystemItem{stexmodel.NewUpLink(dir.ParentDir())}, items...)
	}

	items = filterItems(items, cfg)

	return items
}

// filterItems drops the up link never, drops hidden items when the user has not opted in, and
// applies the regex when there is one. The order of the input slice is preserved.
func filterItems(items []stexmodel.FileSystemItem, cfg config.Config) []stexmodel.FileSystemItem {
	if cfg.ShowHidden && cfg.Filter == nil {
		return items
	}
	filtered := make([]stexmodel.FileSystemItem, 0, len(items))
	for _, item := range items {
		if _, ok := item.(*stexmodel.UpLink); ok {
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

// lessItem is the comparison used by the sort code paths. Items are compared on the chosen
// field, ties broken by UID so the order is stable across runs.
func lessItem(a, b stexmodel.FileSystemItem, by config.SortBy) bool {
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

// sortedSlice returns the values of m sorted by the configured field and order. Ties are
// broken by UID.
func sortedSlice[T stexmodel.FileSystemItem](m map[string]T, cfg config.Config) []stexmodel.FileSystemItem {
	out := make([]stexmodel.FileSystemItem, 0, len(m))
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

// TopChildren returns the n largest direct children of dir, sorted by size descending, with
// the sum of their sizes. n is the maximum number of items to return, values <= 0 return nil.
// Used to populate the "Largest Children" segment of the directory info pane.
func TopChildren(dir *stexmodel.Dir, n int) *core.ChildrenInfo {
	if dir == nil || n <= 0 {
		return nil
	}
	cfg := config.Config{
		SortBy:    config.SortBySize,
		SortOrder: config.Descending,
		Grouping:  config.Mixed,
	}
	items := ComputeItems(dir, cfg)
	filtered := make([]stexmodel.FileSystemItem, 0, len(items))
	for _, item := range items {
		if _, ok := item.(*stexmodel.UpLink); !ok {
			filtered = append(filtered, item)
		}
	}
	if len(filtered) > n {
		filtered = filtered[:n]
	}
	var total stexmodel.Size
	for _, item := range filtered {
		total += item.Size()
	}
	return &core.ChildrenInfo{Items: filtered, TotalSize: total}
}
