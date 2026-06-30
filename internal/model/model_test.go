package model

import (
	"os"
	"path/filepath"
	"testing"
)

// buildTestTree creates a small tree in a temp dir and runs
// BuildTree against it. Returns the root Dir.
func buildTestTree(t *testing.T) *Dir {
	t.Helper()
	tmp := t.TempDir()
	for _, name := range []string{"alpha.txt", "beta.txt", "gamma.txt"} {
		if err := os.WriteFile(filepath.Join(tmp, name), []byte("payload"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(tmp, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "sub", "deep.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	state := &ScanState{}
	BuildTree(tmp, state)
	if state.Result == nil {
		t.Fatal("scan produced no root")
	}
	return state.Result
}

func TestBuildTreePopulatesRoot(t *testing.T) {
	root := buildTestTree(t)
	if root.FullPath() == "" {
		t.Fatal("expected root to have a path")
	}
	if len(root.Files()) != 3 {
		t.Fatalf("expected 3 files at root, got %d", len(root.Files()))
	}
	if _, ok := root.Dirs()["sub"]; !ok {
		t.Fatal("expected sub directory at root")
	}
}

func TestBuildTreeSetsDone(t *testing.T) {
	state := &ScanState{}
	BuildTree(t.TempDir(), state)
	if !state.Done {
		t.Fatal("expected Done to be true after scan")
	}
	if state.Result == nil {
		t.Fatal("expected Result to be set")
	}
}

func TestDirSizeRecursive(t *testing.T) {
	root := buildTestTree(t)
	// three files at root with "payload" (7 bytes each) + one
	// file in sub with "x" (1 byte) = 22 bytes total
	if got := root.Size(); got != Size(22) {
		t.Fatalf("expected size 22, got %d", got)
	}
}

func TestDirSizeCached(t *testing.T) {
	root := buildTestTree(t)
	first := root.Size()
	root.invalidateCaches()
	second := root.Size()
	if first != second {
		t.Fatalf("expected cached and recomputed sizes to match, got %d vs %d", first, second)
	}
}

func TestDirCountRecursive(t *testing.T) {
	root := buildTestTree(t)
	count := root.Count()
	if count.Files != 4 {
		t.Fatalf("expected 4 files total, got %d", count.Files)
	}
	if count.Dirs != 1 {
		t.Fatalf("expected 1 subdirectory, got %d", count.Dirs)
	}
}

func TestDirUIDsAreUnique(t *testing.T) {
	root := buildTestTree(t)
	seen := map[uint64]string{root.UID(): root.Name()}
	for _, file := range root.Files() {
		if existing, ok := seen[file.UID()]; ok {
			t.Fatalf("UID collision between %q and %q", existing, file.Name())
		}
		seen[file.UID()] = file.Name()
	}
	for name, sub := range root.Dirs() {
		if existing, ok := seen[sub.UID()]; ok {
			t.Fatalf("UID collision between %q and %q", existing, name)
		}
		seen[sub.UID()] = name
	}
}

func TestDirRemoveChildFile(t *testing.T) {
	root := buildTestTree(t)
	original := len(root.Files())
	if !root.RemoveChild("alpha.txt") {
		t.Fatal("expected RemoveChild to succeed for an existing file")
	}
	if len(root.Files()) != original-1 {
		t.Fatalf("expected file count to drop by 1, got %d", len(root.Files()))
	}
}

func TestDirRemoveChildDir(t *testing.T) {
	root := buildTestTree(t)
	if !root.RemoveChild("sub") {
		t.Fatal("expected RemoveChild to succeed for an existing directory")
	}
	if _, ok := root.Dirs()["sub"]; ok {
		t.Fatal("expected sub to be removed")
	}
}

func TestDirRemoveChildMissing(t *testing.T) {
	root := buildTestTree(t)
	if root.RemoveChild("nope.txt") {
		t.Fatal("expected RemoveChild to return false for missing child")
	}
}

func TestDirRemoveChildInvalidatesCaches(t *testing.T) {
	root := buildTestTree(t)
	_ = root.Size()
	root.RemoveChild("alpha.txt")
	// After invalidation, Size should recompute and reflect the
	// smaller file count.
	if got := root.Size(); got != Size(15) {
		t.Fatalf("expected size 15 after removing alpha.txt, got %d", got)
	}
}

func TestDirCopyIsDeep(t *testing.T) {
	root := buildTestTree(t)
	clone := root.Copy()
	clone.RemoveChild("alpha.txt")
	if _, ok := root.Files()["alpha.txt"]; !ok {
		t.Fatal("expected the original root to be untouched by clone mutation")
	}
	// Parent pointers in the clone must point at the clone, not
	// the original.
	for _, file := range clone.Files() {
		if file.ParentDir() != clone {
			t.Fatalf("expected clone file parent to be clone, got %p", file.ParentDir())
		}
	}
}

func TestDirCopyIsRecursive(t *testing.T) {
	root := buildTestTree(t)
	clone := root.Copy()
	sub, ok := clone.Dirs()["sub"]
	if !ok {
		t.Fatal("expected sub to be present in clone")
	}
	if sub.ParentDir() != clone {
		t.Fatal("expected clone sub parent to be the clone root")
	}
}

func TestDirSyncDetectsNewFile(t *testing.T) {
	root := buildTestTree(t)
	tmp := root.FullPath()
	if err := os.WriteFile(filepath.Join(tmp, "new.txt"), []byte("z"), 0o644); err != nil {
		t.Fatal(err)
	}
	changed, err := root.Sync()
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("expected Sync to report a change")
	}
	if _, ok := root.Files()["new.txt"]; !ok {
		t.Fatal("expected new.txt to appear after sync")
	}
}

func TestDirSyncDetectsRemovedFile(t *testing.T) {
	root := buildTestTree(t)
	tmp := root.FullPath()
	if err := os.Remove(filepath.Join(tmp, "alpha.txt")); err != nil {
		t.Fatal(err)
	}
	changed, err := root.Sync()
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("expected Sync to report a change")
	}
	if _, ok := root.Files()["alpha.txt"]; ok {
		t.Fatal("expected alpha.txt to be gone after sync")
	}
}

func TestDirSyncReturnsErrorForMissingRoot(t *testing.T) {
	tmp := t.TempDir()
	state := &ScanState{}
	BuildTree(tmp, state)
	root := state.Result
	if err := os.RemoveAll(tmp); err != nil {
		t.Fatal(err)
	}
	_, err := root.Sync()
	if err == nil {
		t.Fatal("expected an error when the root path no longer exists")
	}
}

func TestUpLinkName(t *testing.T) {
	root := buildTestTree(t)
	link := NewUpLink(root)
	if link.Name() != ".." {
		t.Fatalf("expected name \"..\", got %q", link.Name())
	}
}

func TestUpLinkSizeIsZero(t *testing.T) {
	root := buildTestTree(t)
	link := NewUpLink(root)
	if link.Size() != 0 {
		t.Fatalf("expected size 0, got %d", link.Size())
	}
}

func TestUpLinkUIDIsZero(t *testing.T) {
	root := buildTestTree(t)
	link := NewUpLink(root)
	if link.UID() != 0 {
		t.Fatalf("expected UID 0, got %d", link.UID())
	}
}

func TestUpLinkFullPathIsParent(t *testing.T) {
	root := buildTestTree(t)
	link := NewUpLink(root)
	if link.FullPath() != root.FullPath() {
		t.Fatal("expected FullPath to match parent")
	}
}

func TestSizeString(t *testing.T) {
	cases := []struct {
		input Size
		want  string
	}{
		{0, "0 B"},
		{512, "512.00 B"},
		{1024, "1.00 KB"},
		{1024 * 1024, "1.00 MB"},
		{1536, "1.50 KB"},
	}
	for _, c := range cases {
		if got := c.input.String(); got != c.want {
			t.Errorf("Size(%d).String() = %q, want %q", c.input, got, c.want)
		}
	}
}

func TestSizePercentOf(t *testing.T) {
	if got := Size(50).PercentOf(Size(200)); got != 25 {
		t.Errorf("expected 25%%, got %f", got)
	}
	if got := Size(50).PercentOf(Size(0)); got != 0 {
		t.Errorf("expected 0 for zero total, got %f", got)
	}
}
