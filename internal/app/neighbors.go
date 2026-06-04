package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newNeighborsCommand renders the ARP/NDP neighbor table.
func newNeighborsCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "neighbors",
		Aliases: []string{"neighbor", "neigh", "arp"},
		Short:   "Show the ARP/NDP neighbor table",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runNeighbors() },
	}
}

func (o *Options) collectNeighbors(ctx context.Context) ([]model.Neighbor, []model.Warning) {
	return collect.NewNeighbors(o.runner()).Collect(ctx)
}

func (o *Options) runNeighbors() error {
	ctx, cancel := commandContext()
	defer cancel()

	neighbors, warnings := o.collectNeighbors(ctx)
	rep := newReport()
	rep.Neighbors = neighbors
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Neighbors(os.Stdout, neighbors)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderNeighbors(ro render.Options, neighbors []model.Neighbor) {
	ro.Section(os.Stdout, "Neighbors")
	ro.Neighbors(os.Stdout, neighbors)
}
