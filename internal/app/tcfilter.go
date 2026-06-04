package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/bgrewell/ifscope/internal/sysfs"
	"github.com/spf13/cobra"
)

// newFiltersCommand renders the traffic-control filter table.
func newFiltersCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "filters",
		Aliases: []string{"tcfilter", "tc-filters"},
		Short:   "Show traffic-control filters (classifiers)",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runFilters() },
	}
}

func (o *Options) collectFilters(ctx context.Context) ([]model.TCFilter, []model.Warning) {
	return collect.NewTCFilter(o.runner(), sysfs.OS{}).Collect(ctx)
}

func (o *Options) runFilters() error {
	ctx, cancel := commandContext()
	defer cancel()

	filters, warnings := o.collectFilters(ctx)
	rep := newReport()
	rep.TCFilters = filters
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.TCFilter(os.Stdout, filters)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderFilters(ro render.Options, filters []model.TCFilter) {
	ro.Section(os.Stdout, "TC Filters")
	ro.TCFilter(os.Stdout, filters)
}
