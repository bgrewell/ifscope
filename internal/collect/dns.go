package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

const resolvectlCmd = "resolvectl"

// DNS collects resolver state via resolvectl.
type DNS struct {
	Runner run.Runner
}

// NewDNS returns a DNS collector using r.
func NewDNS(r run.Runner) *DNS { return &DNS{Runner: r} }

// Collect runs `resolvectl status` and parses per-link DNS state. A missing
// resolvectl is a non-fatal warning, since DNS is optional context.
func (c *DNS) Collect(ctx context.Context) ([]model.DNS, []model.Warning) {
	stdout, _, err := c.Runner.Run(ctx, resolvectlCmd, "status")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "resolvectl", Message: "resolvectl not found; DNS data unavailable"}}
		}
		return nil, []model.Warning{{Source: "resolvectl", Message: "resolvectl status failed"}}
	}
	return parse.Resolvectl(stdout), nil
}
