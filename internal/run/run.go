// Package run is the single choke point for subprocess execution.
//
// All collectors invoke external tools (ip, ethtool, lspci, resolvectl,
// ovs-vsctl, ping, ...) through a Runner. The real implementation shells out;
// the test implementation returns canned output, so collectors are unit
// testable without touching the host.
package run

import (
	"context"
	"errors"
	"os/exec"
)

// Runner executes a command and returns its stdout, stderr, and error.
type Runner interface {
	Run(ctx context.Context, name string, args ...string) (stdout, stderr []byte, err error)
}

// ErrNotFound indicates the requested executable was not present on PATH. It
// lets collectors distinguish "tool missing" (emit a warning, degrade) from
// "tool failed" (different handling).
var ErrNotFound = errors.New("executable not found")

// IsNotFound reports whether err indicates a missing executable.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) || errors.Is(err, exec.ErrNotFound)
}
