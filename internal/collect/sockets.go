package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

const ssCmd = "ss"

// Sockets collects listening TCP/UDP sockets via ss.
type Sockets struct {
	Runner run.Runner
}

// NewSockets returns a Sockets collector using r.
func NewSockets(r run.Runner) *Sockets { return &Sockets{Runner: r} }

// Collect runs `ss -tulpnH` and parses the listening sockets. Process details
// require privilege; without it they are simply absent. A missing ss is a
// non-fatal warning.
func (c *Sockets) Collect(ctx context.Context) ([]model.Socket, []model.Warning) {
	stdout, _, err := c.Runner.Run(ctx, ssCmd, "-tulpnH")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "ss", Message: "ss not found; socket data unavailable"}}
		}
		return nil, []model.Warning{{Source: "ss", Message: "ss failed"}}
	}
	return parse.SS(stdout), nil
}
