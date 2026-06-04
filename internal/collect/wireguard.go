package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

const wgCmd = "wg"

// WireGuard collects WireGuard interfaces and peers via `wg`.
type WireGuard struct {
	Runner run.Runner
}

// NewWireGuard returns a WireGuard collector using r.
func NewWireGuard(r run.Runner) *WireGuard { return &WireGuard{Runner: r} }

// Collect runs `wg show all dump` and parses it. Private/pre-shared keys are
// discarded by the parser. A missing wg (wireguard-tools not installed) or lack
// of privilege is a non-fatal warning.
func (c *WireGuard) Collect(ctx context.Context) ([]model.WireGuardDevice, []model.Warning) {
	stdout, _, err := c.Runner.Run(ctx, wgCmd, "show", "all", "dump")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "wireguard", Message: "wg not found (wireguard-tools not installed); WireGuard data unavailable"}}
		}
		return nil, []model.Warning{{Source: "wireguard", Message: "wg show failed (root may be required)"}}
	}
	return parse.WireGuard(stdout), nil
}
