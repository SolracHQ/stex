package model

// UpLink is the synthetic ".." entry prepended to a listing when the current directory has a
// parent. It implements FileSystemItem but carries no file data of its own.
type UpLink struct {
	parent *Dir
}

// NewUpLink returns an UpLink pointing at parent. The parent is expected to be non nil when used
// inside a listing, an UpLink with a nil parent is allowed by the type but has no meaning.
func NewUpLink(parent *Dir) *UpLink {
	return &UpLink{parent: parent}
}

// Name returns the literal two character string "..".
func (link *UpLink) Name() string { return ".." }

// Size returns 0. An UpLink is not a filesystem object, it has no size to report.
func (link *UpLink) Size() Size { return 0 }

// Icon returns an empty string, the UpLink row has no icon.
func (link *UpLink) Icon() string { return "" }

// FullPath returns the absolute path of the parent directory, since following ".." leads there.
func (link *UpLink) FullPath() string { return link.parent.fullPath }

// ParentDir returns the directory that ".." navigates up to. Same as FullPath, the parent of
// the UpLink is the directory it navigates to.
func (link *UpLink) ParentDir() *Dir { return link.parent }

// UID returns 0, a reserved sentinel meaning "no real row". The cursor restore logic in the
// explorer skips items whose UID is 0.
func (link *UpLink) UID() uint64 { return 0 }
