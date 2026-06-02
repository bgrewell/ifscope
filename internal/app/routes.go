package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newRoutesCommand renders the routing table.
func newRoutesCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "routes",
		Aliases: []string{"route"},
		Short:   "Show the routing table",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runRoutes() },
	}
}

func (o *Options) collectRoutes(ctx context.Context) ([]model.Route, []model.Warning) {
	return collect.NewRoutes(o.runner()).Collect(ctx)
}

func (o *Options) runRoutes() error {
	ctx, cancel := commandContext()
	defer cancel()

	routes, warnings := o.collectRoutes(ctx)
	rep := newReport()
	rep.Routes = routes
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Routes(os.Stdout, routes)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderRoutes(ro render.Options, routes []model.Route) {
	ro.Section(os.Stdout, "Routes")
	ro.Routes(os.Stdout, routes)
}
