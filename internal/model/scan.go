package model

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// ScanState holds the live progress of an asynchronous directory scan, plus the result once the
// scan completes. A single ScanState is shared between the scanning goroutine and the goroutine
// that drives the Bubble Tea program, all access goes through Mu.
type ScanState struct {
	Mu          sync.Mutex
	CurrentPath string
	TotalItems  int64
	TotalSize   Size
	Warnings    []string
	Done        bool
	Cancelled   bool
	Result      *Dir
	ResultErr   error
}

// BuildTree scans the directory tree rooted at path on the calling goroutine. Progress,
// cancellation, and the final result are communicated through state. The caller is expected to
// start BuildTree in its own goroutine and read state from a Bubble Tea tick.
//
// Cancellation is cooperative, the scan checks state.Cancelled at every directory boundary. The
// function returns immediately when it sees the flag, leaving state.Done false and state.Result
// nil.
//
// On completion state.Done is true and state.Result is the root of the scanned tree.
// state.ResultErr is set only for a fatal error that prevented the scan from producing a result.
func BuildTree(path string, state *ScanState) {
	absPath, _ := filepath.Abs(path)
	root := &Dir{
		uid:      nextUID(),
		name:     absPath,
		fullPath: absPath,
		files:    make(map[string]*File),
		dirs:     make(map[string]*Dir),
	}

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

		for _, child := range entry.dir.dirs {
			stack = append(stack, stackEntry{path: child.fullPath, dir: child})
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

// readDir reads a single directory and populates parent.files and parent.dirs. Non fatal
// errors reading individual entries are silently skipped so a single permission problem does
// not abort the whole scan.
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
			parent.dirs[name] = &Dir{
				uid:      nextUID(),
				name:     name,
				parent:   parent,
				fullPath: fullPath,
				files:    make(map[string]*File),
				dirs:     make(map[string]*Dir),
			}
			items++
		} else {
			fileSize := Size(info.Size())
			parent.files[name] = &File{
				uid:      nextUID(),
				name:     name,
				size:     fileSize,
				parent:   parent,
				fullPath: fullPath,
			}
			items++
			addedSize += fileSize
		}
	}
	return
}
