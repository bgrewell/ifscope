package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

// Queues collects per-interface channel (queue) and ring-buffer settings via
// ethtool -l and -g.
type Queues struct {
	Runner run.Runner
	FS     sysfs.FS
}

// NewQueues returns a Queues collector using r and fs.
func NewQueues(r run.Runner, fs sysfs.FS) *Queues { return &Queues{Runner: r, FS: fs} }

// Collect enumerates interfaces from sysfs and reads channels/rings for each,
// skipping loopback and interfaces that report neither. A missing ethtool is a
// single non-fatal warning.
func (c *Queues) Collect(ctx context.Context) ([]model.Queues, []model.Warning) {
	entries, err := c.FS.ReadDir("/sys/class/net")
	if err != nil {
		return nil, nil
	}

	var out []model.Queues
	for _, e := range entries {
		name := e.Name()
		if name == "lo" {
			continue
		}

		q := model.Queues{Name: name}
		got := false
		if lout, _, lerr := c.Runner.Run(ctx, ethtoolCmd, "-l", name); lerr == nil {
			q.Combined, q.RxChannels, q.TxChannels = parse.EthtoolChannels(lout)
			got = true
		} else if run.IsNotFound(lerr) {
			return nil, []model.Warning{{Source: "ethtool", Message: "ethtool not found; queue data unavailable"}}
		}
		if gout, _, gerr := c.Runner.Run(ctx, ethtoolCmd, "-g", name); gerr == nil {
			q.RxRing, q.TxRing = parse.EthtoolRings(gout)
			got = true
		}
		if got {
			out = append(out, q)
		}
	}
	return out, nil
}
