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

// newClassesCommand renders the traffic-control shaping-class table.
func newClassesCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "classes",
		Aliases: []string{"tcclass", "tc-classes"},
		Short:   "Show traffic-control shaping classes (htb/hfsc rate/ceil)",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runClasses() },
	}
}

func (o *Options) collectClasses(ctx context.Context) ([]model.TCClass, []model.Warning) {
	return collect.NewTCClass(o.runner(), sysfs.OS{}).Collect(ctx)
}

func (o *Options) runClasses() error {
	ctx, cancel := commandContext()
	defer cancel()

	classes, warnings := o.collectClasses(ctx)
	rep := newReport()
	rep.TCClasses = classes
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.TCClass(os.Stdout, classes)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderClasses(ro render.Options, classes []model.TCClass) {
	ro.Section(os.Stdout, "TC Classes")
	ro.TCClass(os.Stdout, classes)
}
