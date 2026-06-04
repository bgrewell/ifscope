package collect

import (
	"context"
	"fmt"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

// Multicast collects IP multicast group memberships via iproute2.
type Multicast struct {
	Runner run.Runner
}

// NewMulticast returns a Multicast collector using r.
func NewMulticast(r run.Runner) *Multicast { return &Multicast{Runner: r} }

// Collect runs `ip -json maddr show` and parses the IP multicast groups.
func (c *Multicast) Collect(ctx context.Context) ([]model.MulticastGroup, []model.Warning) {
	stdout, stderr, err := c.Runner.Run(ctx, ipCmd, "-json", "maddr", "show")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "ip", Message: "ip command not found; multicast data unavailable", Fatal: true}}
		}
		return nil, []model.Warning{{
			Source:  "ip",
			Message: fmt.Sprintf("ip maddr failed: %v: %s", err, strings.TrimSpace(string(stderr))),
		}}
	}
	groups, perr := parse.MulticastGroups(stdout)
	if perr != nil {
		return nil, []model.Warning{{Source: "ip", Message: perr.Error()}}
	}
	return groups, nil
}
