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

// newPTPCommand renders the hardware-timestamping / PTP capability table.
func newPTPCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "ptp",
		Aliases: []string{"timestamping"},
		Short:   "Show hardware timestamping / PTP capabilities",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runPTP() },
	}
}

func (o *Options) collectPTP(ctx context.Context) ([]model.PTP, []model.Warning) {
	return collect.NewPTP(o.runner(), sysfs.OS{}).Collect(ctx)
}

func (o *Options) runPTP() error {
	ctx, cancel := commandContext()
	defer cancel()

	ptps, warnings := o.collectPTP(ctx)
	rep := newReport()
	rep.PTP = ptps
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.PTP(os.Stdout, ptps)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderPTP(ro render.Options, ptps []model.PTP) {
	ro.Section(os.Stdout, "PTP")
	ro.PTP(os.Stdout, ptps)
}
