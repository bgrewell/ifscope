package app

import (
	"os"

	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newVLANsCommand renders only the VLAN table.
func newVLANsCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "vlans",
		Aliases: []string{"vlan"},
		Short:   "Show VLAN interfaces",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runVLANs() },
	}
}

func (o *Options) runVLANs() error {
	ctx, cancel := commandContext()
	defer cancel()

	_, vlans, warnings := o.collectInterfaces(ctx)
	rep := newReport()
	rep.VLANs = vlans
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.VLANs(os.Stdout, vlans)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}
