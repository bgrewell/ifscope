package collect

import (
	"context"
	"fmt"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

// Stats collects per-interface traffic and error counters via iproute2.
type Stats struct {
	Runner run.Runner
}

// NewStats returns a Stats collector using r.
func NewStats(r run.Runner) *Stats { return &Stats{Runner: r} }

// Collect runs `ip -s -j link show` and parses the per-interface counters.
func (c *Stats) Collect(ctx context.Context) ([]model.InterfaceStats, []model.Warning) {
	stdout, stderr, err := c.Runner.Run(ctx, ipCmd, "-s", "-j", "link", "show")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "ip", Message: "ip command not found; statistics unavailable", Fatal: true}}
		}
		return nil, []model.Warning{{
			Source:  "ip",
			Message: fmt.Sprintf("ip -s link failed: %v: %s", err, strings.TrimSpace(string(stderr))),
		}}
	}
	stats, perr := parse.IPLinkStats(stdout)
	if perr != nil {
		return nil, []model.Warning{{Source: "ip", Message: perr.Error()}}
	}
	return stats, nil
}
