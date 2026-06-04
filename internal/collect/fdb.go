package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

const bridgeCmd = "bridge"

// FDB collects bridge forwarding-database entries.
type FDB struct {
	Runner run.Runner
}

// NewFDB returns an FDB collector using r.
func NewFDB(r run.Runner) *FDB { return &FDB{Runner: r} }

// Collect runs `bridge -json fdb show` and parses it. A missing bridge utility
// is a non-fatal warning.
func (c *FDB) Collect(ctx context.Context) ([]model.FDBEntry, []model.Warning) {
	stdout, _, err := c.Runner.Run(ctx, bridgeCmd, "-json", "fdb", "show")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "bridge", Message: "bridge utility not found; FDB data unavailable"}}
		}
		return nil, []model.Warning{{Source: "bridge", Message: "bridge fdb show failed"}}
	}
	entries, perr := parse.BridgeFDB(stdout)
	if perr != nil {
		return nil, []model.Warning{{Source: "bridge", Message: perr.Error()}}
	}
	return entries, nil
}
