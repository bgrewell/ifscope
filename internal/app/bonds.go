package app

import (
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/bgrewell/ifscope/internal/sysfs"
	"github.com/spf13/cobra"
)

// newBondsCommand renders the bonding table.
func newBondsCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "bonds",
		Aliases: []string{"bond"},
		Short:   "Show bonding masters and their members",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runBonds() },
	}
}

// collectBonds gathers bonds from sysfs.
func (o *Options) collectBonds() ([]model.Bond, []model.Warning) {
	return collect.NewBonds(sysfs.OS{}).Collect()
}

func (o *Options) runBonds() error {
	bonds, warnings := o.collectBonds()
	rep := newReport()
	rep.Bonds = bonds
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Bonds(os.Stdout, bonds)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

// renderBonds prints a titled Bonds section.
func renderBonds(ro render.Options, bonds []model.Bond) {
	ro.Section(os.Stdout, "Bonds")
	ro.Bonds(os.Stdout, bonds)
}
