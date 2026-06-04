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
	bvlans, bvw := o.collectBridgeVLANs(ctx)
	warnings = append(warnings, bvw...)
	tunnels, tw := o.collectTunnels(ctx)
	warnings = append(warnings, tw...)
	wg, wgw := o.collectWireGuard(ctx)
	warnings = append(warnings, wgw...)
	neighbors, nw := o.collectNeighbors(ctx)
	warnings = append(warnings, nw...)
	lldp, lw := o.collectLLDP(ctx)
	warnings = append(warnings, lw...)
	fdb, fw := o.collectFDB(ctx)
	warnings = append(warnings, fw...)
	devlink, dlw := o.collectDevlink(ctx)
	warnings = append(warnings, dlw...)
	qdiscs, qw := o.collectQdisc(ctx)
	warnings = append(warnings, qw...)
	classes, clw := o.collectClasses(ctx)
	warnings = append(warnings, clw...)
	filters, flw := o.collectFilters(ctx)
	warnings = append(warnings, flw...)
	offloads, ow := o.collectOffloads(ctx)
	warnings = append(warnings, ow...)
	queues, quw := o.collectQueues(ctx)
	warnings = append(warnings, quw...)
	irqs, iw := o.collectIRQ()
	warnings = append(warnings, iw...)
	ptps, ptw := o.collectPTP(ctx)
	warnings = append(warnings, ptw...)
	mcast, mw := o.collectMulticast(ctx)
	warnings = append(warnings, mw...)
	mdb, mdbw := o.collectMDB(ctx)
	warnings = append(warnings, mdbw...)
	stats, sw := o.collectStats(ctx)
	warnings = append(warnings, sw...)
	sockets, sockw := o.collectSockets(ctx)
	warnings = append(warnings, sockw...)
	netns, nsw := o.collectNetns(ctx)
	warnings = append(warnings, nsw...)
	ovs := o.maybeOVS(ctx, &warnings, ifaces, vlans)

	rep := newReport()
	rep.Interfaces = ifaces
	rep.VLANs = vlans
	rep.Bonds = bonds
	rep.Bridges = bridges
	rep.BridgeVLANs = bvlans
	rep.Tunnels = tunnels
	rep.WireGuard = wg
	rep.PCIe = devices
	rep.Devlink = devlink
	rep.Routes = routes
	rep.Rules = rules
	rep.Neighbors = neighbors
	rep.LLDP = lldp
	rep.FDB = fdb
	rep.DNS = dns
	rep.OVS = ovs
	rep.Qdiscs = qdiscs
	rep.Offloads = offloads
	rep.Queues = queues
	rep.IRQs = irqs
	rep.PTP = ptps
	rep.TCClasses = classes
	rep.TCFilters = filters
	rep.Multicast = mcast
	rep.MDB = mdb
	rep.Netns = netns
	rep.Stats = stats
	rep.Sockets = sockets
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
		if len(bvlans) > 0 {
			fmt.Fprintln(os.Stdout)
			renderBridgeVLANs(ro, bvlans)
		}
		if len(tunnels) > 0 {
			fmt.Fprintln(os.Stdout)
			renderTunnels(ro, tunnels)
		}
		if len(wg) > 0 {
			fmt.Fprintln(os.Stdout)
			renderWireGuard(ro, wg)
		}
		fmt.Fprintln(os.Stdout)
		renderPCIe(ro, devices)
		if len(devlink) > 0 {
			fmt.Fprintln(os.Stdout)
			renderDevlink(ro, devlink)
		}
		fmt.Fprintln(os.Stdout)
		renderRoutes(ro, routes)
		if len(rules) > 0 {
			fmt.Fprintln(os.Stdout)
			renderRules(ro, rules)
		}
		fmt.Fprintln(os.Stdout)
		renderNeighbors(ro, neighbors)
		if len(lldp) > 0 {
			fmt.Fprintln(os.Stdout)
			renderLLDP(ro, lldp)
		}
		if len(fdb) > 0 {
			fmt.Fprintln(os.Stdout)
			renderFDB(ro, fdb)
		}
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
		if len(qdiscs) > 0 {
			fmt.Fprintln(os.Stdout)
			renderQdisc(ro, qdiscs)
		}
		if len(classes) > 0 {
			fmt.Fprintln(os.Stdout)
			renderClasses(ro, classes)
		}
		if len(filters) > 0 {
			fmt.Fprintln(os.Stdout)
			renderFilters(ro, filters)
		}
		if len(offloads) > 0 {
			fmt.Fprintln(os.Stdout)
			renderOffloads(ro, offloads)
		}
		if len(queues) > 0 {
			fmt.Fprintln(os.Stdout)
			renderQueues(ro, queues)
		}
		if len(irqs) > 0 {
			fmt.Fprintln(os.Stdout)
			renderIRQ(ro, irqs)
		}
		if len(ptps) > 0 {
			fmt.Fprintln(os.Stdout)
			renderPTP(ro, ptps)
		}
		if len(mcast) > 0 {
			fmt.Fprintln(os.Stdout)
			renderMulticast(ro, mcast)
		}
		if len(mdb) > 0 {
			fmt.Fprintln(os.Stdout)
			renderMDB(ro, mdb)
		}
		fmt.Fprintln(os.Stdout)
		renderStats(ro, stats)
		fmt.Fprintln(os.Stdout)
		renderSockets(ro, sockets)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}
