package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

// Offloads collects NIC offload feature states (ethtool -k) for each interface.
type Offloads struct {
	Runner run.Runner
	FS     sysfs.FS
}

// NewOffloads returns an Offloads collector using r and fs.
func NewOffloads(r run.Runner, fs sysfs.FS) *Offloads { return &Offloads{Runner: r, FS: fs} }

// Collect enumerates interfaces from sysfs and runs `ethtool -k` for each,
// skipping loopback and interfaces that report no features. A missing ethtool
// yields a single non-fatal warning.
func (c *Offloads) Collect(ctx context.Context) ([]model.Offloads, []model.Warning) {
	entries, err := c.FS.ReadDir("/sys/class/net")
	if err != nil {
		return nil, nil
	}

	var out []model.Offloads
	for _, e := range entries {
		name := e.Name()
		if name == "lo" {
			continue
		}
		stdout, _, rerr := c.Runner.Run(ctx, ethtoolCmd, "-k", name)
		if rerr != nil {
			if run.IsNotFound(rerr) {
				return nil, []model.Warning{{Source: "ethtool", Message: "ethtool not found; offload data unavailable"}}
			}
			continue // interface doesn't support -k
		}
		if features := parse.EthtoolFeatures(stdout); features != nil {
			out = append(out, model.Offloads{Name: name, Features: features})
		}
	}
	return out, nil
}
