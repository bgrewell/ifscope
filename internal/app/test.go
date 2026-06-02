package app

import (
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

// newTestCommand runs basic connectivity tests and binds their flags.
func newTestCommand(o *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run basic connectivity tests",
		Args:  cobra.NoArgs,
		RunE:  func(_ *cobra.Command, _ []string) error { return o.runTest() },
	}

	f := cmd.Flags()
	f.StringVar(&o.PingTarget, "ping-target", collect.DefaultPingTarget, "IP to ping without DNS")
	f.StringVar(&o.DNSTarget, "dns-target", collect.DefaultDNSTarget, "hostname to ping (exercises DNS)")
	f.StringVar(&o.WebTarget, "web-target", collect.DefaultWebTarget, "URL for the HTTPS GET test")
	f.StringVar(&o.ThroughputTarget, "throughput-target", collect.DefaultThroughputTarget, "URL for the throughput test")
	f.BoolVar(&o.Throughput, "throughput", false, "run the download throughput test")
	f.IntVar(&o.Count, "count", 4, "ping packet count")
	f.DurationVar(&o.Timeout, "timeout", 0, "per-test timeout (default 5s)")

	return cmd
}

func (o *Options) connOptions() collect.ConnOptions {
	return collect.ConnOptions{
		PingTarget:       o.PingTarget,
		DNSTarget:        o.DNSTarget,
		WebTarget:        o.WebTarget,
		ThroughputTarget: o.ThroughputTarget,
		Throughput:       o.Throughput,
		Count:            o.Count,
		Timeout:          o.Timeout,
	}
}

func (o *Options) runTest() error {
	ctx, cancel := commandContext()
	defer cancel()

	tests := collect.NewConnectivity(o.runner(), nil).Run(ctx, o.connOptions())
	rep := newReport()
	rep.Tests = tests

	if err := o.emit(rep, func(ro render.Options) {
		ro.Tests(os.Stdout, tests)
	}); err != nil {
		return err
	}
	return exitForFailedTests(tests)
}

// exitForFailedTests returns a non-zero exit error when any test failed.
func exitForFailedTests(tests []model.TestResult) error {
	for _, t := range tests {
		if t.Status == model.StatusFail {
			return &ExitCodeError{Code: ExitTestFailed, Err: errTestFailed}
		}
	}
	return nil
}
