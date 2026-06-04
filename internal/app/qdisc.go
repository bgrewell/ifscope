package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newQdiscCommand renders the per-device queueing-discipline table.
func newQdiscCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "qdisc",
		Aliases: []string{"qdiscs", "shaping"},
		Short:   "Show per-interface queueing disciplines (traffic shaping)",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runQdisc() },
	}
}

func (o *Options) collectQdisc(ctx context.Context) ([]model.Qdisc, []model.Warning) {
	return collect.NewQdisc(o.runner()).Collect(ctx)
}

func (o *Options) runQdisc() error {
	ctx, cancel := commandContext()
	defer cancel()

	qdiscs, warnings := o.collectQdisc(ctx)
	rep := newReport()
	rep.Qdiscs = qdiscs
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Qdisc(os.Stdout, qdiscs)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderQdisc(ro render.Options, qdiscs []model.Qdisc) {
	ro.Section(os.Stdout, "Qdisc")
	ro.Qdisc(os.Stdout, qdiscs)
}
