package collect

import (
	"context"
	"fmt"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

// Tunnels collects overlay/tunnel interfaces (VXLAN, GENEVE, GRE, …).
type Tunnels struct {
	Runner run.Runner
}

// NewTunnels returns a Tunnels collector using r.
func NewTunnels(r run.Runner) *Tunnels { return &Tunnels{Runner: r} }

// Collect runs `ip -d link show` and parses tunnel interfaces from the text
// output. Text is used deliberately: iproute2 emits malformed JSON for vxlan.
func (c *Tunnels) Collect(ctx context.Context) ([]model.Tunnel, []model.Warning) {
	stdout, stderr, err := c.Runner.Run(ctx, ipCmd, "-d", "link", "show")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "ip", Message: "ip command not found; tunnel data unavailable", Fatal: true}}
		}
		return nil, []model.Warning{{
			Source:  "ip",
			Message: fmt.Sprintf("ip -d link failed: %v: %s", err, strings.TrimSpace(string(stderr))),
		}}
	}
	return parse.Tunnels(stdout), nil
}
