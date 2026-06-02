package sysfs

import (
	"io/fs"
	"os"
	"sort"
)

// Fake is an in-memory FS for tests. Files maps path→contents, Links maps a
// symlink path→target, and Dirs maps a directory path→child entry names.
type Fake struct {
	Files map[string]string
	Links map[string]string
	Dirs  map[string][]string
}

// NewFake returns an empty Fake.
func NewFake() *Fake {
	return &Fake{
		Files: map[string]string{},
		Links: map[string]string{},
		Dirs:  map[string][]string{},
	}
}

// ReadFile returns the contents registered for name.
func (f *Fake) ReadFile(name string) ([]byte, error) {
	if v, ok := f.Files[name]; ok {
		return []byte(v), nil
	}
	return nil, &os.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

// Readlink returns the target registered for the symlink name.
func (f *Fake) Readlink(name string) (string, error) {
	if v, ok := f.Links[name]; ok {
		return v, nil
	}
	return "", &os.PathError{Op: "readlink", Path: name, Err: fs.ErrNotExist}
}

// Exists reports whether name is a known file, link, or directory.
func (f *Fake) Exists(name string) bool {
	if _, ok := f.Files[name]; ok {
		return true
	}
	if _, ok := f.Links[name]; ok {
		return true
	}
	_, ok := f.Dirs[name]
	return ok
}

// ReadDir returns the entries registered for the directory name.
func (f *Fake) ReadDir(name string) ([]os.DirEntry, error) {
	names, ok := f.Dirs[name]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	sorted := append([]string(nil), names...)
	sort.Strings(sorted)
	out := make([]os.DirEntry, 0, len(sorted))
	for _, n := range sorted {
		full := name + "/" + n
		_, isLink := f.Links[full]
		_, isDir := f.Dirs[full]
		out = append(out, fakeEntry{name: n, link: isLink, dir: isDir})
	}
	return out, nil
}

// fakeEntry is a minimal os.DirEntry for Fake.ReadDir.
type fakeEntry struct {
	name string
	link bool
	dir  bool
}

func (e fakeEntry) Name() string { return e.name }
func (e fakeEntry) IsDir() bool  { return e.dir }
func (e fakeEntry) Type() fs.FileMode {
	switch {
	case e.link:
		return fs.ModeSymlink
	case e.dir:
		return fs.ModeDir
	default:
		return 0
	}
}
func (e fakeEntry) Info() (fs.FileInfo, error) { return nil, nil }
