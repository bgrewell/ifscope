package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newWireGuardCommand renders the WireGuard interface/peer table.
func newWireGuardCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "wireguard",
		Aliases: []string{"wg"},
		Short:   "Show WireGuard interfaces and peers",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runWireGuard() },
	}
}

func (o *Options) collectWireGuard(ctx context.Context) ([]model.WireGuardDevice, []model.Warning) {
	return collect.NewWireGuard(o.runner()).Collect(ctx)
}

func (o *Options) runWireGuard() error {
	ctx, cancel := commandContext()
	defer cancel()

	devices, warnings := o.collectWireGuard(ctx)
	rep := newReport()
	rep.WireGuard = devices
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.WireGuard(os.Stdout, devices)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderWireGuard(ro render.Options, devices []model.WireGuardDevice) {
	ro.Section(os.Stdout, "WireGuard")
	ro.WireGuard(os.Stdout, devices)
}
