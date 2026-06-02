package app

import (
	"fmt"
	"os"

	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newAllCommand renders every available view. Sections are added as later
// milestones land (routes, DNS, PCIe, OVS, SR-IOV, tests).
func newAllCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Show all available tables",
		Args:  cobra.NoArgs,
		RunE:  func(_ *cobra.Command, _ []string) error { return o.runAll() },
	}
}

func (o *Options) runAll() error {
	ctx, cancel := commandContext()
	defer cancel()

	ifaces, vlans, warnings := o.collectInterfaces(ctx)
	rep := newReport()
	rep.Interfaces = ifaces
	rep.VLANs = vlans
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		renderInterfaces(ro, ifaces)
		fmt.Fprintln(os.Stdout)
		renderVLANs(ro, vlans)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}
