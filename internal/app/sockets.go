package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newSocketsCommand renders the listening-socket table.
func newSocketsCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "sockets",
		Aliases: []string{"listen", "ss"},
		Short:   "Show listening TCP/UDP sockets",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runSockets() },
	}
}

func (o *Options) collectSockets(ctx context.Context) ([]model.Socket, []model.Warning) {
	return collect.NewSockets(o.runner()).Collect(ctx)
}

func (o *Options) runSockets() error {
	ctx, cancel := commandContext()
	defer cancel()

	sockets, warnings := o.collectSockets(ctx)
	rep := newReport()
	rep.Sockets = sockets
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Sockets(os.Stdout, sockets)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderSockets(ro render.Options, sockets []model.Socket) {
	ro.Section(os.Stdout, "Sockets")
	ro.Sockets(os.Stdout, sockets)
}
