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

// newQueuesCommand renders the per-interface channel/ring table.
func newQueuesCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "queues",
		Aliases: []string{"channels"},
		Short:   "Show per-interface queues (channels) and ring sizes",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runQueues() },
	}
}

func (o *Options) collectQueues(ctx context.Context) ([]model.Queues, []model.Warning) {
	return collect.NewQueues(o.runner(), sysfs.OS{}).Collect(ctx)
}

func (o *Options) runQueues() error {
	ctx, cancel := commandContext()
	defer cancel()

	queues, warnings := o.collectQueues(ctx)
	rep := newReport()
	rep.Queues = queues
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Queues(os.Stdout, queues)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderQueues(ro render.Options, queues []model.Queues) {
	ro.Section(os.Stdout, "Queues")
	ro.Queues(os.Stdout, queues)
}
