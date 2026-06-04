package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newMulticastCommand renders the IP multicast group-membership table.
func newMulticastCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "multicast",
		Aliases: []string{"maddr", "mcast"},
		Short:   "Show IP multicast group memberships",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runMulticast() },
	}
}

func (o *Options) collectMulticast(ctx context.Context) ([]model.MulticastGroup, []model.Warning) {
	return collect.NewMulticast(o.runner()).Collect(ctx)
}

func (o *Options) runMulticast() error {
	ctx, cancel := commandContext()
	defer cancel()

	groups, warnings := o.collectMulticast(ctx)
	rep := newReport()
	rep.Multicast = groups
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Multicast(os.Stdout, groups)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderMulticast(ro render.Options, groups []model.MulticastGroup) {
	ro.Section(os.Stdout, "Multicast")
	ro.Multicast(os.Stdout, groups)
}
