package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newDevlinkCommand renders the devlink port table.
func newDevlinkCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "devlink",
		Short: "Show devlink ports (PF/VF flavour, switchdev)",
		Args:  cobra.NoArgs,
		RunE:  func(_ *cobra.Command, _ []string) error { return o.runDevlink() },
	}
}

func (o *Options) collectDevlink(ctx context.Context) ([]model.DevlinkPort, []model.Warning) {
	return collect.NewDevlink(o.runner()).Collect(ctx)
}

func (o *Options) runDevlink() error {
	ctx, cancel := commandContext()
	defer cancel()

	ports, warnings := o.collectDevlink(ctx)
	rep := newReport()
	rep.Devlink = ports
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Devlink(os.Stdout, ports)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderDevlink(ro render.Options, ports []model.DevlinkPort) {
	ro.Section(os.Stdout, "Devlink")
	ro.Devlink(os.Stdout, ports)
}
