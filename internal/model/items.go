// Package model provides the data types that represent a scanned filesystem tree. Dir and File
// hold the tree structure, Size handles human readable byte formatting, and FileSystemItem
// defines the common interface for listing and sorting. The tree is built by BuildTree and
// queried through methods on Dir.
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
