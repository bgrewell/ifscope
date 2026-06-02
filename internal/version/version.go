// Package version exposes build metadata injected at link time.
//
// Values are overridden via -ldflags, e.g.:
//
//	go build -ldflags "-X github.com/bgrewell/ifscope/internal/version.Version=1.2.3 \
//	  -X github.com/bgrewell/ifscope/internal/version.Commit=abc1234 \
//	  -X github.com/bgrewell/ifscope/internal/version.Date=2026-06-02"
package version

import (
	"fmt"
	"runtime"
)

// Build metadata. Defaults apply to non-release builds.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// String returns a one-line, human-readable version summary.
func String() string {
	return fmt.Sprintf("ifscope %s (commit %s, built %s, %s)",
		Version, Commit, Date, runtime.Version())
}
