package collect

import (
	"context"
	"strings"

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
		if cout, _, cerr := c.Runner.Run(ctx, ethtoolCmd, "-c", name); cerr == nil {
			co := parse.EthtoolCoalesce(cout)
			q.RxUsecs, q.TxUsecs = co.RxUsecs, co.TxUsecs
			q.AdaptiveRx, q.AdaptiveTx = co.AdaptiveRx, co.AdaptiveTx
			got = true
		}
		if xout, _, xerr := c.Runner.Run(ctx, ethtoolCmd, "-x", name); xerr == nil {
			q.RSSRings = parse.EthtoolRSSRings(xout)
		}
		q.RPSQueues, q.XPSQueues = c.steeringCounts(name)

		if got {
			out = append(out, q)
		}
	}
	return out, nil
}

// steeringCounts counts the rx queues with a non-zero RPS mask and tx queues
// with a non-zero XPS mask under /sys/class/net/<name>/queues.
func (c *Queues) steeringCounts(name string) (rps, xps int) {
	entries, err := c.FS.ReadDir("/sys/class/net/" + name + "/queues")
	if err != nil {
		return 0, 0
	}
	for _, e := range entries {
		q := e.Name()
		switch {
		case strings.HasPrefix(q, "rx-"):
			if v, err := sysfs.ReadString(c.FS, "/sys/class/net/"+name+"/queues/"+q+"/rps_cpus"); err == nil && maskNonZero(v) {
				rps++
			}
		case strings.HasPrefix(q, "tx-"):
			if v, err := sysfs.ReadString(c.FS, "/sys/class/net/"+name+"/queues/"+q+"/xps_cpus"); err == nil && maskNonZero(v) {
				xps++
			}
		}
	}
	return rps, xps
}

// maskNonZero reports whether a sysfs CPU bitmask (e.g. "00000000,00000000")
// has any bit set.
func maskNonZero(mask string) bool {
	for _, r := range mask {
		if r != '0' && r != ',' && r != ' ' && r != '\n' {
			return true
		}
	}
	return false
}
