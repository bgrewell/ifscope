package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

// MDB collects bridge multicast-database entries.
type MDB struct {
	Runner run.Runner
}

// NewMDB returns an MDB collector using r.
func NewMDB(r run.Runner) *MDB { return &MDB{Runner: r} }

// Collect runs `bridge -json mdb show` and parses it. A missing bridge utility
// is a non-fatal warning.
func (c *MDB) Collect(ctx context.Context) ([]model.MDBEntry, []model.Warning) {
	stdout, _, err := c.Runner.Run(ctx, bridgeCmd, "-json", "mdb", "show")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "bridge", Message: "bridge utility not found; MDB data unavailable"}}
		}
		return nil, []model.Warning{{Source: "bridge", Message: "bridge mdb show failed"}}
	}
	entries, perr := parse.BridgeMDB(stdout)
	if perr != nil {
		return nil, []model.Warning{{Source: "bridge", Message: perr.Error()}}
	}
	return entries, nil
}
