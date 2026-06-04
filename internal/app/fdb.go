package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newFDBCommand renders the bridge forwarding-database table.
func newFDBCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "fdb",
		Short: "Show the bridge forwarding database (MAC table)",
		Args:  cobra.NoArgs,
		RunE:  func(_ *cobra.Command, _ []string) error { return o.runFDB() },
	}
}

func (o *Options) collectFDB(ctx context.Context) ([]model.FDBEntry, []model.Warning) {
	return collect.NewFDB(o.runner()).Collect(ctx)
}

func (o *Options) runFDB() error {
	ctx, cancel := commandContext()
	defer cancel()

	entries, warnings := o.collectFDB(ctx)
	rep := newReport()
	rep.FDB = entries
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.FDB(os.Stdout, entries)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderFDB(ro render.Options, entries []model.FDBEntry) {
	ro.Section(os.Stdout, "FDB")
	ro.FDB(os.Stdout, entries)
}
