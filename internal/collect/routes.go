package collect

import (
	"context"
	"fmt"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

// Routes collects the kernel routing table via iproute2.
type Routes struct {
	Runner run.Runner
}

// NewRoutes returns a Routes collector using r.
func NewRoutes(r run.Runner) *Routes { return &Routes{Runner: r} }

// Collect runs `ip -detail -json route show table all` so routes from every
// table (local, default, and custom tables used by policy routing) are
// included, then parses them. iproute2 omits the table name for main-table
// routes, so an empty table is normalized to "main".
func (c *Routes) Collect(ctx context.Context) ([]model.Route, []model.Warning) {
	stdout, stderr, err := c.Runner.Run(ctx, ipCmd, "-detail", "-json", "route", "show", "table", "all")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "ip", Message: "ip command not found; route data unavailable", Fatal: true}}
		}
		return nil, []model.Warning{{
			Source:  "ip",
			Message: fmt.Sprintf("ip route failed: %v: %s", err, strings.TrimSpace(string(stderr))),
		}}
	}
	routes, perr := parse.IPRoutes(stdout)
	if perr != nil {
		return nil, []model.Warning{{Source: "ip", Message: perr.Error()}}
	}
	for i := range routes {
		if routes[i].Table == "" {
			routes[i].Table = "main"
		}
	}
	return routes, nil
}
