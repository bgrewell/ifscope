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
	rules, rulw := o.collectRules(ctx)
	warnings = append(warnings, rulw...)
	dns, dw := o.collectDNS(ctx)
	warnings = append(warnings, dw...)
	bonds, bw := o.collectBonds()
	warnings = append(warnings, bw...)
	bridges, brw := o.collectBridges()
	warnings = append(warnings, brw...)
	neighbors, nw := o.collectNeighbors(ctx)
	warnings = append(warnings, nw...)
	stats, sw := o.collectStats(ctx)
	warnings = append(warnings, sw...)
	netns, nsw := o.collectNetns(ctx)
	warnings = append(warnings, nsw...)
	ovs := o.maybeOVS(ctx, &warnings, ifaces, vlans)

	rep := newReport()
	rep.Interfaces = ifaces
	rep.VLANs = vlans
	rep.Bonds = bonds
	rep.Bridges = bridges
	rep.PCIe = devices
	rep.Routes = routes
	rep.Rules = rules
	rep.Neighbors = neighbors
	rep.DNS = dns
	rep.OVS = ovs
	rep.Netns = netns
	rep.Stats = stats
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
		fmt.Fprintln(os.Stdout)
		renderPCIe(ro, devices)
		fmt.Fprintln(os.Stdout)
		renderRoutes(ro, routes)
		if len(rules) > 0 {
			fmt.Fprintln(os.Stdout)
			renderRules(ro, rules)
		}
		fmt.Fprintln(os.Stdout)
		renderNeighbors(ro, neighbors)
		fmt.Fprintln(os.Stdout)
		renderDNS(ro, dns)
		if ovs != nil {
			fmt.Fprintln(os.Stdout)
			renderOVS(ro, ovs)
		}
		if len(netns) > 0 {
			fmt.Fprintln(os.Stdout)
			renderNetns(ro, netns)
		}
		fmt.Fprintln(os.Stdout)
		renderStats(ro, stats)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}
