package app

import (
	"context"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/correlate"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newStatsCommand renders per-interface traffic and error counters.
func newStatsCommand(o *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stats",
		Aliases: []string{"statistics", "counters"},
		Short:   "Show per-interface traffic and error counters",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runStats() },
	}
	cmd.Flags().DurationVar(&o.StatsRate, "rate", 0, "sample for this duration and show per-second rates (100ms-30s)")
	cmd.Flags().StringVar(&o.StatsSort, "sort", "name", "sort by: name|rx|tx|errors|drops|utilization")
	cmd.Flags().IntVar(&o.StatsTop, "top", 0, "show only the first N interfaces after sorting (0 = all)")
	return cmd
}

func (o *Options) collectStats(ctx context.Context) ([]model.InterfaceStats, []model.Warning) {
	return collect.NewStats(o.runner()).Collect(ctx)
}

func (o *Options) runStats() error {
	if err := o.validateStatsOptions(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout+o.StatsRate)
	defer cancel()

	stats, warnings := o.collectStats(ctx)
	sampledAt := time.Now()
	if o.StatsRate > 0 {
		first, firstWarnings := stats, warnings
		timer := time.NewTimer(o.StatsRate)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
		stats, warnings = o.collectStats(ctx)
		warnings = append(firstWarnings, warnings...)
		stats = correlate.StatsRates(first, stats, time.Since(sampledAt))
	} else if o.Watch > 0 {
		if o.statsPrevious != nil {
			stats = correlate.StatsRates(o.statsPrevious, stats, sampledAt.Sub(o.statsPreviousAt))
		}
		o.statsPrevious = append([]model.InterfaceStats(nil), stats...)
		for i := range o.statsPrevious {
			o.statsPrevious[i].Rates = nil
			o.statsPrevious[i].RateStatus = ""
		}
		o.statsPreviousAt = sampledAt
	}
	correlate.SortStats(stats, o.StatsSort)
	if o.StatsTop > 0 && len(stats) > o.StatsTop {
		stats = stats[:o.StatsTop]
	}
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

func (o *Options) validateStatsOptions() error {
	if o.StatsRate > 0 && (o.StatsRate < 100*time.Millisecond || o.StatsRate > 30*time.Second) {
		return fmt.Errorf("--rate must be between 100ms and 30s")
	}
	if o.StatsRate > 0 && o.Watch > 0 {
		return fmt.Errorf("--rate and --watch cannot be used together")
	}
	if !slices.Contains([]string{"name", "rx", "tx", "errors", "drops", "utilization"}, o.StatsSort) {
		return fmt.Errorf("invalid --sort %q: use name, rx, tx, errors, drops, or utilization", o.StatsSort)
	}
	if o.StatsTop < 0 {
		return fmt.Errorf("--top cannot be negative")
	}
	if o.StatsSort != "name" && o.StatsRate == 0 && o.Watch == 0 {
		return fmt.Errorf("--sort %s requires --rate or --watch", o.StatsSort)
	}
	return nil
}

func renderStats(ro render.Options, stats []model.InterfaceStats) {
	ro.Section(os.Stdout, "Statistics")
	ro.Stats(os.Stdout, stats)
}
