package app

import (
	"context"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newDNSCommand renders the DNS resolver table.
func newDNSCommand(o *Options) *cobra.Command {
	return &cobra.Command{
		Use:   "dns",
		Short: "Show DNS resolver state",
		Args:  cobra.NoArgs,
		RunE:  func(_ *cobra.Command, _ []string) error { return o.runDNS() },
	}
}

func (o *Options) collectDNS(ctx context.Context) ([]model.DNS, []model.Warning) {
	return collect.NewDNS(o.runner()).Collect(ctx)
}

func (o *Options) runDNS() error {
	ctx, cancel := commandContext()
	defer cancel()

	dns, warnings := o.collectDNS(ctx)
	rep := newReport()
	rep.DNS = dns
	rep.Warnings = warnings

	if err := o.emit(rep, func(ro render.Options) {
		ro.DNS(os.Stdout, dns)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}

func renderDNS(ro render.Options, dns []model.DNS) {
	ro.Section(os.Stdout, "DNS")
	ro.DNS(os.Stdout, dns)
}
