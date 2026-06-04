package app

import (
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/bgrewell/ifscope/internal/sysfs"
	"github.com/spf13/cobra"
)

// newIRQCommand renders the NIC interrupt-affinity table.
func newIRQCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "irq",
		Aliases: []string{"irqs", "affinity"},
		Short:   "Show NIC interrupt CPU affinity",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runIRQ() },
	}
}

func (o *Options) collectIRQ() ([]model.IRQ, []model.Warning) {
	return collect.NewIRQ(sysfs.OS{}).Collect()
}

func (o *Options) runIRQ() error {
	irqs, warnings := o.collectIRQ()
	rep := newReport()
	rep.IRQs = irqs
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.IRQ(os.Stdout, irqs)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderIRQ(ro render.Options, irqs []model.IRQ) {
	ro.Section(os.Stdout, "IRQ affinity")
	ro.IRQ(os.Stdout, irqs)
}
