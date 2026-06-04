package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

const devlinkCmd = "devlink"

// Devlink collects devlink port information.
type Devlink struct {
	Runner run.Runner
}

// NewDevlink returns a Devlink collector using r.
func NewDevlink(r run.Runner) *Devlink { return &Devlink{Runner: r} }

// Collect runs `devlink -j port show` and parses it. A missing devlink or a
// host without devlink-capable devices is a non-fatal warning.
func (c *Devlink) Collect(ctx context.Context) ([]model.DevlinkPort, []model.Warning) {
	stdout, _, err := c.Runner.Run(ctx, devlinkCmd, "-j", "port", "show")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "devlink", Message: "devlink not found; devlink data unavailable"}}
		}
		return nil, []model.Warning{{Source: "devlink", Message: "devlink port show failed"}}
	}
	ports, perr := parse.DevlinkPorts(stdout)
	if perr != nil {
		return nil, []model.Warning{{Source: "devlink", Message: perr.Error()}}
	}
	return ports, nil
}
