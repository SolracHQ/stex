package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func defaultPathForTest(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	path, err := DefaultPath()
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func TestLoadMissingFileReturnsErrNotFound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	_, err := Load(path)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSaveThenLoadRoundTrip(t *testing.T) {
	path := defaultPathForTest(t)
	want := Config{
		SortBy:     SortByName,
		SortOrder:  Ascending,
		Grouping:   DirsFirst,
		ShowIcons:  true,
		ShowHidden: true,
		LiveFilter: false,
	}
	if err := want.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got != want {
		t.Fatalf("round trip mismatch:\nwant %+v\ngot  %+v", want, got)
	}
}

func TestSaveSkipsFilter(t *testing.T) {
	path := defaultPathForTest(t)
	c := DefaultConfig()
	if err := c.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) == 0 || body[0] != '{' {
		t.Fatalf("expected JSON object, got %q", body)
	}
	if contains(body, []byte("Filter")) {
		t.Fatalf("file should not contain Filter key: %s", body)
	}
}

func TestSaveCreatesParentDir(t *testing.T) {
	_ = defaultPathForTest(t)
	if err := DefaultConfig().Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	path, _ := DefaultPath()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file at %s: %v", path, err)
	}
}

func TestLoadCorruptFileReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte("not json {{{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}

func contains(haystack, needle []byte) bool {
	if len(needle) == 0 {
		return true
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if string(haystack[i:i+len(needle)]) == string(needle) {
			return true
		}
	}
	return false
}
