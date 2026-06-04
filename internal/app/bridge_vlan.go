package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newBridgeVLANsCommand renders the bridge VLAN-filtering table.
func newBridgeVLANsCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "bridge-vlans",
		Aliases: []string{"bvlan"},
		Short:   "Show bridge VLAN-filtering entries",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runBridgeVLANs() },
	}
}

func (o *Options) collectBridgeVLANs(ctx context.Context) ([]model.BridgeVLAN, []model.Warning) {
	return collect.NewBridgeVLANs(o.runner()).Collect(ctx)
}

func (o *Options) runBridgeVLANs() error {
	ctx, cancel := commandContext()
	defer cancel()

	vlans, warnings := o.collectBridgeVLANs(ctx)
	rep := newReport()
	rep.BridgeVLANs = vlans
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.BridgeVLANs(os.Stdout, vlans)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderBridgeVLANs(ro render.Options, vlans []model.BridgeVLAN) {
	ro.Section(os.Stdout, "Bridge VLANs")
	ro.BridgeVLANs(os.Stdout, vlans)
}
