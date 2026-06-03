package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newRulesCommand renders the routing policy rule table.
func newRulesCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:     "rules",
		Aliases: []string{"rule"},
		Short:   "Show routing policy rules (source-based routing)",
		Args:    cobra.NoArgs,
		RunE:    func(_ *cobra.Command, _ []string) error { return o.runRules() },
	}
}

func (o *Options) collectRules(ctx context.Context) ([]model.Rule, []model.Warning) {
	return collect.NewRules(o.runner()).Collect(ctx)
}

func (o *Options) runRules() error {
	ctx, cancel := commandContext()
	defer cancel()

	rules, warnings := o.collectRules(ctx)
	rep := newReport()
	rep.Rules = rules
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.Rules(os.Stdout, rules)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

// renderRules prints a titled Rules section.
func renderRules(ro render.Options, rules []model.Rule) {
	ro.Section(os.Stdout, "Rules")
	ro.Rules(os.Stdout, rules)
}
