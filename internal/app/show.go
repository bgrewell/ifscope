package app

import (
	"fmt"
	"os"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newShowCommand renders the default view: interfaces and VLANs.
func newShowCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show interfaces and VLANs (default view)",
		Args:  cobra.NoArgs,
		RunE:  func(_ *cobra.Command, _ []string) error { return o.runShow() },
	}
}

// runShow collects interfaces+VLANs and emits both sections.
func (o *Options) runShow() error {
	ctx, cancel := commandContext()
	defer cancel()

	ifaces, vlans, warnings := o.collectInterfaces(ctx)
	bonds, bw := o.collectBonds()
	warnings = append(warnings, bw...)
	bridges, brw := o.collectBridges()
	warnings = append(warnings, brw...)
	ovs := o.maybeOVS(ctx, &warnings, ifaces, vlans)

	rep := newReport()
	rep.Interfaces = ifaces
	rep.VLANs = vlans
	rep.Bonds = bonds
	rep.Bridges = bridges
	rep.OVS = ovs
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		renderInterfaces(ro, ifaces)
		fmt.Fprintln(os.Stdout)
		renderVLANs(ro, vlans)
		if len(bonds) > 0 {
			fmt.Fprintln(os.Stdout)
			renderBonds(ro, bonds)
		}
		if len(bridges) > 0 {
			fmt.Fprintln(os.Stdout)
			renderBridges(ro, bridges)
		}
		if ovs != nil {
			fmt.Fprintln(os.Stdout)
			renderOVS(ro, ovs)
		}
	}); err != nil {
		return err
	}
	return exitForWarnings(rep.Warnings)
}

func renderInterfaces(ro render.Options, ifaces []model.Interface) {
	ro.Section(os.Stdout, "Interfaces")
	ro.Interfaces(os.Stdout, ifaces)
}

func renderVLANs(ro render.Options, vlans []model.Interface) {
	ro.Section(os.Stdout, "VLANs")
	ro.VLANs(os.Stdout, vlans)
}

func renderPCIe(ro render.Options, devices []model.PCIDevice) {
	ro.Section(os.Stdout, "PCIe")
	ro.PCIe(os.Stdout, devices)
}
