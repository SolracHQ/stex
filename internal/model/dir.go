package model

import (
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
)

// uidCounter is a process wide counter used to mint stable UIDs for every Dir and File in the
// model. UIDs only need to be unique within a single run, they are not persisted.
var uidCounter atomic.Uint64

func nextUID() uint64 {
	return uidCounter.Add(1)
}

// Dir represents a directory in the scanned tree. Children are stored as maps keyed by base name,
// which gives O(1) name lookup and preserves the order returned by the filesystem reader. A Dir
// is always owned by exactly one parent except the root, whose parent is nil.
type Dir struct {
	uid             uint64
	name            string
	parent          *Dir
	files           map[string]*File
	dirs            map[string]*Dir
	fullPath        string
	lastSelectedUID uint64
	size            Size
	count           Count
}

// UID returns the stable identifier of the directory. UIDs are unique across the lifetime of a
// process and survive resyncs, which lets callers restore a cursor to the same logical row
// across mode transitions and tree refreshes.
func (dir *Dir) UID() uint64 { return dir.uid }

// Name returns the base name of the directory with a trailing slash, the canonical visual form
// for a directory row.
func (dir *Dir) Name() string { return dir.name + "/" }

// FullPath returns the absolute path of the directory on disk.
func (dir *Dir) FullPath() string { return dir.fullPath }

// ParentDir returns the parent directory, or nil when this Dir is the root of the scanned tree.
func (dir *Dir) ParentDir() *Dir { return dir.parent }

// Icon returns a short emoji used to label the row in the explorer. The value is constant for
// the directory kind.
func (dir *Dir) Icon() string { return "📁" }

// LastSelectedUID returns the UID of the row that was highlighted the last time the cursor was
// inside this directory. 0 means no row was selected, which is the case for a freshly scanned
// directory.
func (dir *Dir) LastSelectedUID() uint64 { return dir.lastSelectedUID }

// SetLastSelectedUID records the UID of the currently highlighted row, so the explorer can
// restore the cursor when navigating back into this directory later.
func (dir *Dir) SetLastSelectedUID(uid uint64) { dir.lastSelectedUID = uid }

// Files returns the directory's file children as a map keyed by base name. The returned map is
// the live internal map, callers that need to mutate it must hold a reference to the parent Dir
// and use SetFiles to replace it.
func (dir *Dir) Files() map[string]*File { return dir.files }

// Dirs returns the directory's subdirectory children as a map keyed by base name. The returned
// map is the live internal map, callers that need to mutate it must hold a reference to the
// parent Dir and use SetDirs to replace it.
func (dir *Dir) Dirs() map[string]*Dir { return dir.dirs }

// SetFiles replaces the file children map. The new map becomes the live internal map. The
// caller is responsible for keeping each file's parent pointer consistent.
func (dir *Dir) SetFiles(files map[string]*File) { dir.files = files }

// SetDirs replaces the subdirectory children map. The new map becomes the live internal map.
// The caller is responsible for keeping each subdirectory's parent pointer consistent.
func (dir *Dir) SetDirs(dirs map[string]*Dir) { dir.dirs = dirs }

// RemoveChild removes the direct child with the given base name. name is the basename as stored
// in the maps, not a full path. Returns true if a child was removed, false if no child had that
// name. The cached size and count are invalidated on success.
func (dir *Dir) RemoveChild(name string) bool {
	if _, ok := dir.files[name]; ok {
		delete(dir.files, name)
		dir.invalidateCaches()
		return true
	}
	if _, ok := dir.dirs[name]; ok {
		delete(dir.dirs, name)
		dir.invalidateCaches()
		return true
	}
	return false
}

// Size returns the total size in bytes of all files under the directory, recursively. The
// result is cached after the first call, callers that mutate the tree must call InvalidateUp
// to force recomputation.
func (dir *Dir) Size() Size {
	if dir.size > 0 {
		return dir.size
	}
	var total Size
	for _, file := range dir.files {
		total += file.size
	}
	for _, sub := range dir.dirs {
		total += sub.Size()
	}
	dir.size = total
	return total
}

// invalidateCaches clears the cached size and count so the next call to Size or Count
// recomputes them.
func (dir *Dir) invalidateCaches() {
	dir.size = 0
	dir.count = Count{}
}

// InvalidateUp walks the chain from dir to the root and clears the cached size and count on
// every ancestor. Call after a subtree mutation so all affected levels recompute on demand.
func (dir *Dir) InvalidateUp() {
	for cur := dir; cur != nil; cur = cur.parent {
		cur.invalidateCaches()
	}
}

// Count holds the total number of files and subdirectories under a directory, recursively.
type Count struct {
	Files int
	Dirs  int
}

// Count returns the total number of files and subdirectories under the directory, recursively.
// The result is cached after the first call, callers that mutate the tree must call
// InvalidateUp to force recomputation.
func (dir *Dir) Count() Count {
	if dir.count.Files != 0 || dir.count.Dirs != 0 {
		return dir.count
	}
	dir.count = Count{Files: len(dir.files), Dirs: len(dir.dirs)}
	for _, sub := range dir.dirs {
		subCount := sub.Count()
		dir.count.Files += subCount.Files
		dir.count.Dirs += subCount.Dirs
	}
	return dir.count
}

// Copy returns a deep copy of the directory subtree with every parent pointer rewired to the
// copy. The returned tree is safe to mutate without affecting the original. Cached size and
// count are copied as is and will be recomputed on demand if invalidated.
func (dir *Dir) Copy() *Dir {
	clone := &Dir{
		uid:             dir.uid,
		name:            dir.name,
		fullPath:        dir.fullPath,
		lastSelectedUID: dir.lastSelectedUID,
		size:            dir.size,
		count:           dir.count,
		files:           make(map[string]*File, len(dir.files)),
		dirs:            make(map[string]*Dir, len(dir.dirs)),
	}
	for name, file := range dir.files {
		clone.files[name] = &File{
			uid:      file.uid,
			name:     file.name,
			size:     file.size,
			parent:   clone,
			fullPath: file.fullPath,
		}
	}
	for name, sub := range dir.dirs {
		clone.dirs[name] = sub.Copy()
		clone.dirs[name].parent = clone
	}
	return clone
}

// Sync reconciles the directory subtree with the filesystem. It returns true if anything in
// the subtree changed, and an error if the root of this subtree no longer exists on disk. A
// nil error means the path was readable and was either left as is, or updated to match the
// disk contents.
//
// On a successful sync the cached size and count are invalidated at every level that changed.
func (dir *Dir) Sync() (bool, error) {
	if _, err := os.Stat(dir.fullPath); err != nil {
		if dir.parent == nil {
			return false, fmt.Errorf("root directory %q was removed", dir.fullPath)
		}
		delete(dir.parent.dirs, dir.name)
		dir.parent.invalidateCaches()
		return true, nil
	}

	entries, readErr := os.ReadDir(dir.fullPath)
	if readErr != nil {
		return false, nil
	}

	current := make(map[string]bool, len(entries))
	for _, entry := range entries {
		current[entry.Name()] = entry.IsDir()
	}

	changed := false

	for name, file := range dir.files {
		if !current[name] {
			delete(dir.files, name)
			changed = true
			continue
		}
		info, err := os.Stat(file.fullPath)
		if err == nil && Size(info.Size()) != file.size {
			file.size = Size(info.Size())
			changed = true
		}
	}

	for name := range dir.dirs {
		if !current[name] {
			delete(dir.dirs, name)
			changed = true
		}
	}

	for name, isDir := range current {
		_, hasFile := dir.files[name]
		_, hasDir := dir.dirs[name]
		if hasFile || hasDir {
			continue
		}
		fullPath := filepath.Join(dir.fullPath, name)
		if isDir {
			dir.dirs[name] = &Dir{
				uid:      nextUID(),
				name:     name,
				parent:   dir,
				fullPath: fullPath,
				files:    make(map[string]*File),
				dirs:     make(map[string]*Dir),
			}
		} else {
			info, err := os.Stat(fullPath)
			if err == nil {
				dir.files[name] = &File{
					uid:      nextUID(),
					name:     name,
					size:     Size(info.Size()),
					parent:   dir,
					fullPath: fullPath,
				}
			}
		}
		changed = true
	}

	for _, sub := range dir.dirs {
		subChanged, err := sub.Sync()
		if err != nil {
			return false, err
		}
		if subChanged {
			changed = true
		}
	}

	if changed {
		dir.invalidateCaches()
	}
	return changed, nil
}
