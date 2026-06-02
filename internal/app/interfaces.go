package app

import (
	"fmt"
	"os"

	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newInterfacesCommand renders only the interface table.
func newInterfacesCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "interfaces",
		Aliases: []string{"if", "iface"},
		Short:   "Show network interfaces",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runInterfaces() },
	}
}

func (o *Options) runInterfaces() error {
	ctx, cancel := commandContext()
	defer cancel()

	ifaces, _, warnings := o.collectInterfaces(ctx)
	ovs := o.maybeOVS(ctx, &warnings, ifaces)

	rep := newReport()
	rep.Interfaces = ifaces
	rep.OVS = ovs
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Interfaces(os.Stdout, ifaces)
		if ovs != nil {
			fmt.Fprintln(os.Stdout)
			renderOVS(ro, ovs)
		}
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}
