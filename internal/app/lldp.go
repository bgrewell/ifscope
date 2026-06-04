package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newLLDPCommand renders the LLDP neighbor table.
func newLLDPCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "lldp",
		Short: "Show LLDP link-layer neighbors",
		Args:  cobra.NoArgs,
		RunE:  func(_ *cobra.Command, _ []string) error { return o.runLLDP() },
	}
}

func (o *Options) collectLLDP(ctx context.Context) ([]model.LLDPNeighbor, []model.Warning) {
	return collect.NewLLDP(o.runner()).Collect(ctx)
}

func (o *Options) runLLDP() error {
	ctx, cancel := commandContext()
	defer cancel()

	neighbors, warnings := o.collectLLDP(ctx)
	rep := newReport()
	rep.LLDP = neighbors
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.LLDP(os.Stdout, neighbors)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderLLDP(ro render.Options, neighbors []model.LLDPNeighbor) {
	ro.Section(os.Stdout, "LLDP")
	ro.LLDP(os.Stdout, neighbors)
}
