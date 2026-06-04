package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newTunnelsCommand renders the tunnel/overlay interface table.
func newTunnelsCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "tunnels",
		Aliases: []string{"tunnel"},
		Short:   "Show overlay/tunnel interfaces (VXLAN, GENEVE, GRE, ...)",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runTunnels() },
	}
}

func (o *Options) collectTunnels(ctx context.Context) ([]model.Tunnel, []model.Warning) {
	return collect.NewTunnels(o.runner()).Collect(ctx)
}

func (o *Options) runTunnels() error {
	ctx, cancel := commandContext()
	defer cancel()

	tunnels, warnings := o.collectTunnels(ctx)
	rep := newReport()
	rep.Tunnels = tunnels
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Tunnels(os.Stdout, tunnels)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderTunnels(ro render.Options, tunnels []model.Tunnel) {
	ro.Section(os.Stdout, "Tunnels")
	ro.Tunnels(os.Stdout, tunnels)
}
