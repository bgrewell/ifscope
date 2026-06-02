package app

import (
	"fmt"
	"os"

	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// compatFlags holds the legacy netcheck short flags preserved for muscle
// memory. They are hidden from help; the subcommands are the documented UI.
//
// -j/-u/-s already exist as short forms of --json/--up/--summary. The view
// selectors below let `ifscope -I`, `ifscope -IP`, etc. mirror netcheck's
// flag-combining behavior when no subcommand is given.
type compatFlags struct {
	interfaces bool
	vlans      bool
	dns        bool
	routes     bool
	pcie       bool
	test       bool
}

// bindCompat registers the hidden legacy view-selector flags on the root.
func (c *compatFlags) bindCompat(cmd *cobra.Command) {
	f := cmd.Flags()
	f.BoolVarP(&c.interfaces, "compat-interfaces", "I", false, "(compat) show interfaces")
	f.BoolVarP(&c.vlans, "compat-vlans", "V", false, "(compat) show VLANs")
	f.BoolVarP(&c.dns, "compat-dns", "D", false, "(compat) show DNS")
	f.BoolVarP(&c.routes, "compat-routes", "R", false, "(compat) show routes")
	f.BoolVarP(&c.pcie, "compat-pcie", "P", false, "(compat) show PCIe")
	f.BoolVarP(&c.test, "compat-test", "t", false, "(compat) run connectivity tests")
	for _, name := range []string{
		"compat-interfaces", "compat-vlans", "compat-dns",
		"compat-routes", "compat-pcie", "compat-test",
	} {
		_ = f.MarkHidden(name)
	}
}

// any reports whether at least one legacy view selector is set.
func (c compatFlags) any() bool {
	return c.interfaces || c.vlans || c.dns || c.routes || c.pcie || c.test
}

// runCompat renders the sections selected by legacy view flags. Sections are
// wired in as their milestones land; until then a selector for an unimplemented
// view is a no-op.
func (o *Options) runCompat(c compatFlags) error {
	ctx, cancel := commandContext()
	defer cancel()

	rep := newReport()
	var sections []func(render.Options)

	if c.interfaces || c.vlans || c.pcie {
		ifaces, vlans, devices, w := o.gather(ctx, c.pcie)
		rep.Warnings = append(rep.Warnings, w...)
		if c.interfaces {
			rep.Interfaces = ifaces
			sections = append(sections, func(ro render.Options) { renderInterfaces(ro, ifaces) })
		}
		if c.vlans {
			rep.VLANs = vlans
			sections = append(sections, func(ro render.Options) { renderVLANs(ro, vlans) })
		}
		if c.pcie {
			rep.PCIe = devices
			sections = append(sections, func(ro render.Options) { renderPCIe(ro, devices) })
		}
	}

	if c.routes {
		routes, w := o.collectRoutes(ctx)
		rep.Warnings = append(rep.Warnings, w...)
		rep.Routes = routes
		sections = append(sections, func(ro render.Options) { renderRoutes(ro, routes) })
	}
	if c.dns {
		dns, w := o.collectDNS(ctx)
		rep.Warnings = append(rep.Warnings, w...)
		rep.DNS = dns
		sections = append(sections, func(ro render.Options) { renderDNS(ro, dns) })
	}

	if err := o.emit(rep, func(ro render.Options) {
		for i, fn := range sections {
			if i > 0 {
				fmt.Fprintln(os.Stdout)
			}
			fn(ro)
		}
	}); err != nil {
		return err
	}
	return exitForWarnings(rep.Warnings)
}
