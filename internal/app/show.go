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

func renderInterfaces(ro render.Options, ifaces []model.Interface) {
	ro.Section(os.Stdout, "Interfaces")
	ro.Interfaces(os.Stdout, ifaces)
}

func renderVLANs(ro render.Options, vlans []model.Interface) {
	ro.Section(os.Stdout, "VLANs")
	ro.VLANs(os.Stdout, vlans)
}
