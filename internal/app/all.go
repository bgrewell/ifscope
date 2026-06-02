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

	ifaces, vlans, devices, warnings := o.gather(ctx, true)
	routes, rw := o.collectRoutes(ctx)
	warnings = append(warnings, rw...)
	dns, dw := o.collectDNS(ctx)
	warnings = append(warnings, dw...)
	ovs := o.maybeOVS(ctx, &warnings, ifaces, vlans)

	rep := newReport()
	rep.Interfaces = ifaces
	rep.VLANs = vlans
	rep.PCIe = devices
	rep.Routes = routes
	rep.DNS = dns
	rep.OVS = ovs
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		renderInterfaces(ro, ifaces)
		fmt.Fprintln(os.Stdout)
		renderVLANs(ro, vlans)
		fmt.Fprintln(os.Stdout)
		renderPCIe(ro, devices)
		fmt.Fprintln(os.Stdout)
		renderRoutes(ro, routes)
		fmt.Fprintln(os.Stdout)
		renderDNS(ro, dns)
		if ovs != nil {
			fmt.Fprintln(os.Stdout)
			renderOVS(ro, ovs)
		}
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}
