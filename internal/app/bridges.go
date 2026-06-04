package app

import (
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/bgrewell/ifscope/internal/sysfs"
	"github.com/spf13/cobra"
)

// newBridgesCommand renders the Linux bridge table.
func newBridgesCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "bridges",
		Aliases: []string{"bridge", "br"},
		Short:   "Show Linux bridges and their member ports",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runBridges() },
	}
}

// collectBridges gathers Linux bridges from sysfs.
func (o *Options) collectBridges() ([]model.Bridge, []model.Warning) {
	return collect.NewBridges(sysfs.OS{}).Collect()
}

func (o *Options) runBridges() error {
	bridges, warnings := o.collectBridges()
	rep := newReport()
	rep.Bridges = bridges
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Bridges(os.Stdout, bridges)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

// renderBridges prints a titled Bridges section.
func renderBridges(ro render.Options, bridges []model.Bridge) {
	ro.Section(os.Stdout, "Bridges")
	ro.Bridges(os.Stdout, bridges)
}
