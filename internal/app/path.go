package app

import (
	"fmt"
	"os"

	"github.com/bgrewell/ifscope/internal/collect"
	"github.com/bgrewell/ifscope/internal/render"
	"github.com/spf13/cobra"
)

type pathFlags struct {
	family   string
	source   string
	outIface string
	protocol string
	port     int
}

func newPathCommand(o *Options) *cobra.Command {
	flags := &pathFlags{}
	cmd := &cobra.Command{
		Use:   "path DESTINATION",
		Short: "Explain the route and local path to a destination",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return o.runPath(args[0], *flags)
		},
	}
	f := cmd.Flags()
	f.StringVar(&flags.family, "family", "any", "address family: any, 4, or 6")
	f.StringVar(&flags.source, "source", "", "source address for route selection")
	f.StringVar(&flags.outIface, "out-interface", "", "constrain route lookup to an egress interface")
	f.StringVar(&flags.protocol, "protocol", "", "transport protocol for policy lookup: tcp or udp")
	f.IntVar(&flags.port, "port", 0, "destination port for policy lookup")
	return cmd
}

func (o *Options) runPath(destination string, flags pathFlags) error {
	if flags.family != "any" && flags.family != "4" && flags.family != "6" {
		return fmt.Errorf("--family must be any, 4, or 6")
	}
	if flags.protocol != "" && flags.protocol != "tcp" && flags.protocol != "udp" {
		return fmt.Errorf("--protocol must be tcp or udp")
	}
	if flags.port < 0 || flags.port > 65535 {
		return fmt.Errorf("--port must be between 1 and 65535")
	}
	if flags.port != 0 && flags.protocol == "" {
		return fmt.Errorf("--port requires --protocol")
	}

	ctx, cancel := commandContext()
	defer cancel()
	path, warnings := collect.NewPath(o.runner()).Collect(ctx, destination, collect.PathOptions{
		Family: flags.family, Source: flags.source, OutIface: flags.outIface,
		Protocol: flags.protocol, Port: flags.port,
	})
	rep := newReport()
	rep.Path = &path
	rep.Warnings = warnings
	if err := o.emit(rep, func(ro render.Options) {
		ro.Path(os.Stdout, path)
	}); err != nil {
		return err
	}
	return exitForWarnings(warnings)
}
