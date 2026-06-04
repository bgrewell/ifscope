package collect

import (
	"context"
	"fmt"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

// Neighbors collects the ARP/NDP neighbor table via iproute2.
type Neighbors struct {
	Runner run.Runner
}

// NewNeighbors returns a Neighbors collector using r.
func NewNeighbors(r run.Runner) *Neighbors { return &Neighbors{Runner: r} }

// Collect runs `ip -json neigh show` and parses it.
func (c *Neighbors) Collect(ctx context.Context) ([]model.Neighbor, []model.Warning) {
	stdout, stderr, err := c.Runner.Run(ctx, ipCmd, "-json", "neigh", "show")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "ip", Message: "ip command not found; neighbor data unavailable", Fatal: true}}
		}
		return nil, []model.Warning{{
			Source:  "ip",
			Message: fmt.Sprintf("ip neigh failed: %v: %s", err, strings.TrimSpace(string(stderr))),
		}}
	}
	neighbors, perr := parse.IPNeighbors(stdout)
	if perr != nil {
		return nil, []model.Warning{{Source: "ip", Message: perr.Error()}}
	}
	return neighbors, nil
}
