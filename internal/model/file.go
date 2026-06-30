package model

// File represents a regular file in the scanned tree. A File is always owned by exactly one
// parent Dir, the parent's files map holds the reference keyed by the file's base name.
type File struct {
	uid      uint64
	name     string
	size     Size
	parent   *Dir
	fullPath string
}

// UID returns the stable identifier of the file. UIDs are unique across the lifetime of a
// process and survive resyncs, which lets callers restore a cursor to the same logical row
// across mode transitions and tree refreshes.
func (file *File) UID() uint64 { return file.uid }

// Name returns the base name of the file without any path separator, for example "report.txt".
func (file *File) Name() string { return file.name }

// Size returns the file's size in bytes as recorded at the most recent scan or resync. The
// value is not live, it does not reflect changes to the file on disk until the next sync.
func (file *File) Size() Size { return file.size }

// FullPath returns the absolute path of the file on disk.
func (file *File) FullPath() string { return file.fullPath }

// ParentDir returns the directory that owns this file. Never nil since every file in the model
// has a parent.
func (file *File) ParentDir() *Dir { return file.parent }

// Icon returns a short emoji used to label the row in the explorer. The value is constant for
// the file kind.
func (file *File) Icon() string { return "📄" }
