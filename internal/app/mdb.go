package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newMDBCommand renders the bridge multicast-database table.
func newMDBCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "mdb",
		Short: "Show the bridge multicast database",
		Args:  cobra.NoArgs,
		RunE:  func(_ *cobra.Command, _ []string) error { return o.runMDB() },
	}
}

func (o *Options) collectMDB(ctx context.Context) ([]model.MDBEntry, []model.Warning) {
	return collect.NewMDB(o.runner()).Collect(ctx)
}

func (o *Options) runMDB() error {
	ctx, cancel := commandContext()
	defer cancel()

	entries, warnings := o.collectMDB(ctx)
	rep := newReport()
	rep.MDB = entries
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.MDB(os.Stdout, entries)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderMDB(ro render.Options, entries []model.MDBEntry) {
	ro.Section(os.Stdout, "MDB")
	ro.MDB(os.Stdout, entries)
}
