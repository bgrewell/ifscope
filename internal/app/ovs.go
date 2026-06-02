package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/correlate"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newOVSCommand renders the Open vSwitch topology.
func newOVSCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "ovs",
		Short: "Show Open vSwitch bridges, ports, and interfaces",
		Args:  cobra.NoArgs,
		RunE:  func(_ *cobra.Command, _ []string) error { return o.runOVS() },
	}
}

func (o *Options) collectOVS(ctx context.Context) (*model.OVS, []model.Warning) {
	return collect.NewOVS(o.runner(), o.NoSudo).Collect(ctx)
}

// ovsEnabled reports whether OVS data should be gathered for a view: when the
// user passes --ovs (and not --no-ovs).
func (o *Options) ovsEnabled() bool { return o.OVS && !o.NoOVS }

// maybeOVS collects and annotates OVS membership when --ovs is enabled,
// appending any warnings. It returns nil when OVS is not requested.
func (o *Options) maybeOVS(ctx context.Context, warnings *[]model.Warning, groups ...[]model.Interface) *model.OVS {
	if !o.ovsEnabled() {
		return nil
	}
	ovs, w := o.collectOVS(ctx)
	*warnings = append(*warnings, w...)
	correlate.AnnotateOVS(ovs, groups...)
	return ovs
}

func (o *Options) runOVS() error {
	ctx, cancel := commandContext()
	defer cancel()

	ovs, warnings := o.collectOVS(ctx)
	rep := newReport()
	rep.OVS = ovs
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.OVS(os.Stdout, ovs)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderOVS(ro render.Options, ovs *model.OVS) {
	ro.Section(os.Stdout, "OVS")
	ro.OVS(os.Stdout, ovs)
}
