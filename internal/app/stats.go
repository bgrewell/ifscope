package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newStatsCommand renders per-interface traffic and error counters.
func newStatsCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "stats",
		Aliases: []string{"statistics", "counters"},
		Short:   "Show per-interface traffic and error counters",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runStats() },
	}
}

func (o *Options) collectStats(ctx context.Context) ([]model.InterfaceStats, []model.Warning) {
	return collect.NewStats(o.runner()).Collect(ctx)
}

func (o *Options) runStats() error {
	ctx, cancel := commandContext()
	defer cancel()

	stats, warnings := o.collectStats(ctx)
	rep := newReport()
	rep.Stats = stats
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Stats(os.Stdout, stats)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderStats(ro render.Options, stats []model.InterfaceStats) {
	ro.Section(os.Stdout, "Statistics")
	ro.Stats(os.Stdout, stats)
}
