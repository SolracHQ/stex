package model

// FileSystemItem is the read only interface satisfied by any value that can appear in a
// directory listing: a file, a subdirectory, or the ".." up link entry. Modes that need to render
// or sort items take []FileSystemItem rather than concrete types.
type FileSystemItem interface {
	Name() string
	Size() Size
	Icon() string
	FullPath() string
	ParentDir() *Dir
	UID() uint64
}
