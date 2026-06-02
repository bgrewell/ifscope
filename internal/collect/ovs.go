package collect

import (
	"context"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

const ovsVsctlCmd = "ovs-vsctl"

// OVS collects Open vSwitch topology via ovs-vsctl. Reads are attempted
// unprivileged first and retried with `sudo -n` on failure, unless NoSudo is
// set. sudo is invoked non-interactively so a password-required host fails
// fast with a warning rather than hanging.
type OVS struct {
	Runner run.Runner
	NoSudo bool
}

// NewOVS returns an OVS collector using r.
func NewOVS(r run.Runner, noSudo bool) *OVS { return &OVS{Runner: r, NoSudo: noSudo} }

// Collect returns the OVS model, or nil with a warning when ovs-vsctl is
// absent or access is denied.
func (c *OVS) Collect(ctx context.Context) (*model.OVS, []model.Warning) {
	bridgeData, w := c.list(ctx, "Bridge", parse.OVSBridgeColumns)
	if w != nil {
		return nil, []model.Warning{*w}
	}
	portData, w := c.list(ctx, "Port", parse.OVSPortColumns)
	if w != nil {
		return nil, []model.Warning{*w}
	}
	ifaceData, w := c.list(ctx, "Interface", parse.OVSInterfaceColumns)
	if w != nil {
		return nil, []model.Warning{*w}
	}

	ovs, err := parse.OVS(bridgeData, portData, ifaceData)
	if err != nil {
		return nil, []model.Warning{{Source: "ovs", Message: err.Error()}}
	}
	return &ovs, nil
}

// list runs `ovs-vsctl --format=json --columns=<cols> list <table>`, retrying
// once via sudo on failure. It returns a warning instead of data when the tool
// is missing or remains inaccessible.
func (c *OVS) list(ctx context.Context, table, cols string) ([]byte, *model.Warning) {
	args := []string{"--format=json", "--columns=" + cols, "list", table}

	out, _, err := c.Runner.Run(ctx, ovsVsctlCmd, args...)
	if err == nil {
		return out, nil
	}
	if run.IsNotFound(err) {
		return nil, &model.Warning{Source: "ovs-vsctl", Message: "ovs-vsctl not found; OVS data unavailable"}
	}
	if c.NoSudo {
		return nil, &model.Warning{Source: "ovs", Message: "ovs-vsctl access denied and --no-sudo set; run as root for OVS data"}
	}

	sudoArgs := append([]string{"-n", ovsVsctlCmd}, args...)
	out, _, serr := c.Runner.Run(ctx, "sudo", sudoArgs...)
	if serr == nil {
		return out, nil
	}
	return nil, &model.Warning{Source: "ovs", Message: "ovs-vsctl requires privileges; run as root or grant sudo for OVS data"}
}
