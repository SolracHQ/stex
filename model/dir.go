package model

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// File represents a regular file in the file tree.
type File struct {
	Name     string
	Size     Size
	Parent   *Dir
	FullPath string
}

// Dir represents a directory in the file tree.
type Dir struct {
	Name       string
	Parent     *Dir // nil for root
	Files      []File
	Dirs       []*Dir
	FullPath   string
	LastCursor int
	size       Size
	count      Count
}

// ScanState holds live progress data for the async directory scan.
type ScanState struct {
	Mu          sync.Mutex
	CurrentPath string   // path currently being scanned
	TotalItems  int64    // cumulative file and directory count
	TotalSize   Size     // cumulative file size
	Warnings    []string // non-fatal scan errors
	Done        bool     // set when scanning completes
	Cancelled   bool     // set when user aborts
	Result      *Dir     // root directory of the scanned tree (set when Done)
	ResultErr   error    // fatal scan error (set when Done)
}

// BuildTree scans the directory tree rooted at path. Progress, cancellation,
// and errors are communicated through state.
func BuildTree(path string, state *ScanState) {
	absPath, _ := filepath.Abs(path)
	root := &Dir{Name: absPath, FullPath: absPath}

	state.Mu.Lock()
	state.CurrentPath = absPath
	state.Mu.Unlock()

	type stackEntry struct {
		path string
		dir  *Dir
	}
	stack := []stackEntry{{path: absPath, dir: root}}

	var (
		totalItems int64
		totalSize  Size
	)

	for len(stack) > 0 {
		state.Mu.Lock()
		if state.Cancelled {
			state.Mu.Unlock()
			return
		}
		state.Mu.Unlock()

		entry := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		state.Mu.Lock()
		state.CurrentPath = entry.path
		state.Mu.Unlock()

		items, size, err := readDir(entry.path, entry.dir)
		if err != nil {
			state.Mu.Lock()
			state.Warnings = append(state.Warnings, fmt.Sprintf("error scanning %s", entry.path))
			state.Mu.Unlock()
		}

		for _, child := range entry.dir.Dirs {
			stack = append(stack, stackEntry{path: child.FullPath, dir: child})
		}
		totalItems += items
		totalSize += size

		state.Mu.Lock()
		state.TotalItems = totalItems
		state.TotalSize = totalSize
		state.Mu.Unlock()
	}

	state.Mu.Lock()
	state.Result = root
	state.Done = true
	state.TotalItems = totalItems
	state.TotalSize = totalSize
	state.Mu.Unlock()
}

// readDir reads a single directory, populating parent.Files and parent.Dirs.
func readDir(path string, parent *Dir) (items int64, addedSize Size, err error) {
	entries, readErr := os.ReadDir(path)
	if readErr != nil {
		return 0, 0, readErr
	}
	for _, dirent := range entries {
		name := dirent.Name()
		fullPath := filepath.Join(path, name)

		info, statErr := dirent.Info()
		if statErr != nil {
			continue
		}

		if dirent.IsDir() {
			child := &Dir{
				Name:     name + "/",
				Parent:   parent,
				FullPath: fullPath,
			}
			parent.Dirs = append(parent.Dirs, child)
			items++
		} else {
			fileSize := Size(info.Size())
			parent.Files = append(parent.Files, File{
				Name:     name,
				Size:     fileSize,
				Parent:   parent,
				FullPath: fullPath,
			})
			items++
			addedSize += fileSize
		}
	}
	return
}

// Size returns the total file size under d.
func (d *Dir) Size() Size {
	if d.size > 0 {
		return d.size
	}
	var total Size
	for _, f := range d.Files {
		total += f.Size
	}
	for _, subDir := range d.Dirs {
		total += subDir.Size()
	}
	d.size = total
	return total
}

// Count holds file and directory totals.
type Count struct {
	Files int
	Dirs  int
}

// Count returns total files and directories under d, cached after first call.
func (d *Dir) Count() Count {
	if d.count.Files != 0 || d.count.Dirs != 0 {
		return d.count
	}
	d.count = Count{Files: len(d.Files), Dirs: len(d.Dirs)}
	for _, subDir := range d.Dirs {
		subCount := subDir.Count()
		d.count.Files += subCount.Files
		d.count.Dirs += subCount.Dirs
	}
	return d.count
}
