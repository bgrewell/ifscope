package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

// PTP collects hardware-timestamping / PTP capabilities via ethtool -T for
// physical interfaces.
type PTP struct {
	Runner run.Runner
	FS     sysfs.FS
}

// NewPTP returns a PTP collector using r and fs.
func NewPTP(r run.Runner, fs sysfs.FS) *PTP { return &PTP{Runner: r, FS: fs} }

// Collect runs `ethtool -T` for each device-backed (physical) interface. A
// missing ethtool is a single non-fatal warning.
func (c *PTP) Collect(ctx context.Context) ([]model.PTP, []model.Warning) {
	entries, err := c.FS.ReadDir("/sys/class/net")
	if err != nil {
		return nil, nil
	}

	var out []model.PTP
	for _, e := range entries {
		name := e.Name()
		// PTP/hw-timestamping is a property of real NICs.
		if !c.FS.Exists("/sys/class/net/" + name + "/device") {
			continue
		}
		stdout, _, rerr := c.Runner.Run(ctx, ethtoolCmd, "-T", name)
		if rerr != nil {
			if run.IsNotFound(rerr) {
				return nil, []model.Warning{{Source: "ethtool", Message: "ethtool not found; PTP data unavailable"}}
			}
			continue
		}
		p := parse.EthtoolTimestamp(stdout)
		p.Name = name
		out = append(out, p)
	}
	return out, nil
}
