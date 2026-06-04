package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

const tcCmd = "tc"

// Qdisc collects queueing disciplines via tc. Only root qdiscs are returned —
// the effective per-device discipline — to avoid the noise of per-queue child
// qdiscs (e.g. the dozens of pfifo_fast under an mq root).
type Qdisc struct {
	Runner run.Runner
}

// NewQdisc returns a Qdisc collector using r.
func NewQdisc(r run.Runner) *Qdisc { return &Qdisc{Runner: r} }

// Collect runs `tc -json qdisc show` and returns the root qdisc per device.
func (c *Qdisc) Collect(ctx context.Context) ([]model.Qdisc, []model.Warning) {
	stdout, _, err := c.Runner.Run(ctx, tcCmd, "-json", "qdisc", "show")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "tc", Message: "tc not found; qdisc data unavailable"}}
		}
		return nil, []model.Warning{{Source: "tc", Message: "tc qdisc show failed"}}
	}
	all, perr := parse.Qdiscs(stdout)
	if perr != nil {
		return nil, []model.Warning{{Source: "tc", Message: perr.Error()}}
	}
	roots := make([]model.Qdisc, 0, len(all))
	for _, q := range all {
		if q.Root {
			roots = append(roots, q)
		}
	}
	return roots, nil
}
