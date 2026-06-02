// Package sysfs abstracts reads from /sys and /proc behind an interface so
// that sysfs-driven logic (PCIe and SR-IOV inspection) is unit testable
// against fixture trees instead of the live host.
package sysfs

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FS provides read-only access to a pseudo-filesystem (sysfs/procfs).
//
// Paths are absolute as seen on the host (e.g. "/sys/class/net/eth0/address").
// Implementations may transparently re-root them for testing.
type FS interface {
	ReadFile(name string) ([]byte, error)
	ReadDir(name string) ([]os.DirEntry, error)
	Readlink(name string) (string, error)
	Exists(name string) bool
}

// OS is the production FS rooted at Root (defaults to "/").
type OS struct {
	// Root is prepended to every path; empty means "/".
	Root string
}

func (o OS) path(name string) string {
	if o.Root == "" {
		return name
	}
	return filepath.Join(o.Root, name)
}

// ReadFile reads the named file.
func (o OS) ReadFile(name string) ([]byte, error) { return os.ReadFile(o.path(name)) }

// ReadDir lists the named directory.
func (o OS) ReadDir(name string) ([]os.DirEntry, error) { return os.ReadDir(o.path(name)) }

// Readlink returns the target of the named symbolic link.
func (o OS) Readlink(name string) (string, error) { return os.Readlink(o.path(name)) }

// Exists reports whether the named path exists.
func (o OS) Exists(name string) bool {
	_, err := os.Lstat(o.path(name))
	return err == nil
}

// ReadString reads a file and returns its trimmed string contents. Sysfs
// attributes conventionally carry a trailing newline.
func ReadString(fs FS, name string) (string, error) {
	b, err := fs.ReadFile(name)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

// ReadInt reads a file and parses its trimmed contents as a base-10 integer.
func ReadInt(fs FS, name string) (int, error) {
	s, err := ReadString(fs, name)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}
