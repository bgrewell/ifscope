package app

import (
	"os"

	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newSRIOVCommand renders the SR-IOV PF/VF table.
func newSRIOVCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "sriov",
		Short: "Show SR-IOV physical and virtual function state",
		Args:  cobra.NoArgs,
		RunE:  func(_ *cobra.Command, _ []string) error { return o.runSRIOV() },
	}
}

func (o *Options) runSRIOV() error {
	ctx, cancel := commandContext()
	defer cancel()

	// SR-IOV needs the full (pre-filter) interface set to enumerate PFs, but
	// the displayed rows honor --pf/--vf/--interface filters via gather().
	ifaces, _, warnings := o.collectInterfaces(ctx)
	rep := newReport()
	rep.Interfaces = ifaces
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.SRIOV(os.Stdout, ifaces)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}
