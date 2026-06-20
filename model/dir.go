package model

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// File represents a regular file in the file tree.
type File struct {
	Name     string // basename of the file
	Size     Size   // size in bytes
	Parent   *Dir   // parent directory
	FullPath string // absolute path
}

// Dir represents a directory in the file tree.
type Dir struct {
	Name     string // basename with trailing "/" for children, full path for root
	Parent   *Dir   // parent directory, nil for root
	Files    []File
	Dirs     []*Dir
	FullPath string  // absolute path
	size     Size    // cached total size, computed lazily by Size()
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

		entries, err := os.ReadDir(entry.path)
		if err != nil {
			state.Mu.Lock()
			state.Warnings = append(state.Warnings, fmt.Sprintf("Warning: error scanning %s", entry.path))
			state.Mu.Unlock()
			continue
		}

		for _, dirent := range entries {
			name := dirent.Name()
			fullPath := filepath.Join(entry.path, name)

			info, err := dirent.Info()
			if err != nil {
				continue
			}

			if dirent.IsDir() {
				child := &Dir{
					Name:     name + "/",
					Parent:   entry.dir,
					FullPath: fullPath,
				}
				entry.dir.Dirs = append(entry.dir.Dirs, child)
				stack = append(stack, stackEntry{path: fullPath, dir: child})
				totalItems++
			} else {
				fileSize := Size(info.Size())
				entry.dir.Files = append(entry.dir.Files, File{
					Name:     name,
					Size:     fileSize,
					Parent:   entry.dir,
					FullPath: fullPath,
				})
				totalItems++
				totalSize += fileSize
			}
		}

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

// Size returns the total file size under d, cached after first call.
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

// Count returns total files and directories under d.
func (d *Dir) Count() Count {
	count := Count{Files: len(d.Files), Dirs: len(d.Dirs)}
	for _, subDir := range d.Dirs {
		subCount := subDir.Count()
		count.Files += subCount.Files
		count.Dirs += subCount.Dirs
	}
	return count
}
