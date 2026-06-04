package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

// BridgeVLANs collects bridge VLAN-filtering entries.
type BridgeVLANs struct {
	Runner run.Runner
}

// NewBridgeVLANs returns a BridgeVLANs collector using r.
func NewBridgeVLANs(r run.Runner) *BridgeVLANs { return &BridgeVLANs{Runner: r} }

// Collect runs `bridge -json vlan show` and parses it. A missing bridge utility
// is a non-fatal warning.
func (c *BridgeVLANs) Collect(ctx context.Context) ([]model.BridgeVLAN, []model.Warning) {
	stdout, _, err := c.Runner.Run(ctx, bridgeCmd, "-json", "vlan", "show")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "bridge", Message: "bridge utility not found; VLAN-filter data unavailable"}}
		}
		return nil, []model.Warning{{Source: "bridge", Message: "bridge vlan show failed"}}
	}
	vlans, perr := parse.BridgeVLANs(stdout)
	if perr != nil {
		return nil, []model.Warning{{Source: "bridge", Message: perr.Error()}}
	}
	return vlans, nil
}
