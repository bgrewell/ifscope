// Package testutil provides shared helpers for locating repository fixtures.
package testutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// repoRoot walks up from this source file until it finds go.mod.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("testutil: cannot determine caller")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("testutil: go.mod not found above test file")
		}
		dir = parent
	}
}

// Fixture reads testdata/<rel> from the repository root, failing the test on
// error. The single fixtures tree is shared across packages.
func Fixture(t *testing.T, rel string) []byte {
	t.Helper()
	path := filepath.Join(repoRoot(t), "testdata", filepath.FromSlash(rel))
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %s: %v", rel, err)
	}
	return b
}
