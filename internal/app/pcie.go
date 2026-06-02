package app

import (
	"os"

	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newPCIeCommand renders the PCIe network-device table.
func newPCIeCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "pcie",
		Aliases: []string{"pci"},
		Short:   "Show PCIe network devices",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runPCIe() },
	}
}

func (o *Options) runPCIe() error {
	ctx, cancel := commandContext()
	defer cancel()

	_, _, devices, warnings := o.gather(ctx, true)
	rep := newReport()
	rep.PCIe = devices
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.PCIe(os.Stdout, devices)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}
