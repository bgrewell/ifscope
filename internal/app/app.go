package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/correlate"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/version"
)

// defaultTimeout bounds a single command's collection work.
const defaultTimeout = 30 * time.Second

// runner returns the production command runner.
func (o *Options) runner() run.Runner { return run.Exec{} }

// commandContext returns a context bounded by the default timeout.
func commandContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultTimeout)
}

// renderOptions maps global flags to render options bound to stdout.
func (o *Options) renderOptions() render.Options {
	return render.Options{
		Summary:   o.Summary,
		Barebones: o.Barebones,
		Color:     render.NewColor(o.Color, o.NoColor, os.Stdout),
	}
}

// filter maps global filter flags to a correlate.Filter.
func (o *Options) filter() correlate.Filter {
	return correlate.Filter{
		Up:       o.Up,
		Name:     o.Interface,
		Driver:   o.Driver,
		State:    o.State,
		Physical: o.Physical,
		Virtual:  o.Virtual,
		VF:       o.VF,
		PF:       o.PF,
	}
}

// newReport returns a report seeded with schema version and host identity.
func newReport() *model.Report {
	host, _ := os.Hostname()
	return &model.Report{
		Version: version.Version,
		Host:    model.Host{Hostname: host},
	}
}

// collectInterfaces gathers interfaces and VLANs, applying filters and sorting.
func (o *Options) collectInterfaces(ctx context.Context) (ifaces, vlans []model.Interface, warnings []model.Warning) {
	all, w := collect.NewInterfaces(o.runner()).Collect(ctx)
	warnings = append(warnings, w...)

	ifaces, vlans = correlate.Partition(all)
	f := o.filter()
	ifaces = f.Apply(ifaces)
	vlans = f.Apply(vlans)
	correlate.SortInterfaces(ifaces)
	correlate.SortVLANs(vlans)
	return ifaces, vlans, warnings
}

// emit writes the report: JSON to stdout in JSON mode, otherwise it invokes
// tables to render human output and then reports warnings to stderr.
func (o *Options) emit(rep *model.Report, tables func(ro render.Options)) error {
	if o.JSON {
		return render.JSON(os.Stdout, rep, o.Pretty)
	}
	tables(o.renderOptions())
	o.reportWarnings(rep.Warnings)
	return nil
}

// exitForWarnings returns a dependency-missing exit error if any warning is
// fatal, so callers can surface a non-zero exit code after emitting output.
func exitForWarnings(warnings []model.Warning) error {
	for _, w := range warnings {
		if w.Fatal {
			return &ExitCodeError{Code: ExitDepMissing, Err: fmt.Errorf("required data could not be collected")}
		}
	}
	return nil
}

// reportWarnings prints warnings to stderr when requested or when any is fatal.
func (o *Options) reportWarnings(warnings []model.Warning) {
	if len(warnings) == 0 {
		return
	}
	show := o.Warnings
	for _, w := range warnings {
		if w.Fatal {
			show = true
		}
	}
	if !show {
		return
	}
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "warning [%s]: %s\n", w.Source, w.Message)
	}
}
