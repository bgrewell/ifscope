package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newNetnsCommand renders the network namespace table.
func newNetnsCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "netns",
		Aliases: []string{"namespaces", "ns"},
		Short:   "Show network namespaces",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runNetns() },
	}
}

func (o *Options) collectNetns(ctx context.Context) ([]model.Netns, []model.Warning) {
	return collect.NewNetns(o.runner()).Collect(ctx)
}

func (o *Options) runNetns() error {
	ctx, cancel := commandContext()
	defer cancel()

	namespaces, warnings := o.collectNetns(ctx)
	rep := newReport()
	rep.Netns = namespaces
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Netns(os.Stdout, namespaces)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderNetns(ro render.Options, namespaces []model.Netns) {
	ro.Section(os.Stdout, "Namespaces")
	ro.Netns(os.Stdout, namespaces)
}
