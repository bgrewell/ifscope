package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

const lldpcliCmd = "lldpcli"

// LLDP collects link-layer (LLDP) neighbors via lldpd's lldpcli.
type LLDP struct {
	Runner run.Runner
}

// NewLLDP returns an LLDP collector using r.
func NewLLDP(r run.Runner) *LLDP { return &LLDP{Runner: r} }

// Collect runs `lldpcli show neighbors -f json0` and parses it. The json0
// format is used because it wraps every field in arrays, avoiding the
// object-vs-array ambiguity of `-f json`. A missing lldpcli (lldpd not
// installed) is a non-fatal warning.
func (c *LLDP) Collect(ctx context.Context) ([]model.LLDPNeighbor, []model.Warning) {
	stdout, _, err := c.Runner.Run(ctx, lldpcliCmd, "show", "neighbors", "-f", "json0")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "lldp", Message: "lldpcli not found (lldpd not installed); LLDP data unavailable"}}
		}
		return nil, []model.Warning{{Source: "lldp", Message: "lldpcli show neighbors failed"}}
	}
	neighbors, perr := parse.LLDPNeighbors(stdout)
	if perr != nil {
		return nil, []model.Warning{{Source: "lldp", Message: perr.Error()}}
	}
	return neighbors, nil
}
