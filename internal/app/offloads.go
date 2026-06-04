package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/bgrewell/ifscope/internal/sysfs"
	"github.com/spf13/cobra"
)

// newOffloadsCommand renders the NIC offload-feature table.
func newOffloadsCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "offloads",
		Aliases: []string{"offload"},
		Short:   "Show NIC offload features (checksum, TSO, GSO, GRO, ...)",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runOffloads() },
	}
}

func (o *Options) collectOffloads(ctx context.Context) ([]model.Offloads, []model.Warning) {
	return collect.NewOffloads(o.runner(), sysfs.OS{}).Collect(ctx)
}

func (o *Options) runOffloads() error {
	ctx, cancel := commandContext()
	defer cancel()

	offloads, warnings := o.collectOffloads(ctx)
	rep := newReport()
	rep.Offloads = offloads
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Offloads(os.Stdout, offloads)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderOffloads(ro render.Options, offloads []model.Offloads) {
	ro.Section(os.Stdout, "Offloads")
	ro.Offloads(os.Stdout, offloads)
}
