// Package app wires the cobra command tree, flag handling, and orchestration
// that turns CLI invocations into collected, rendered reports.
package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/bgrewell/ifscope/internal/version"
	"github.com/spf13/cobra"
)

const longDescription = `ifscope inspects Linux host network state in one place.

It correlates interfaces, VLANs, routes, DNS, PCIe NIC mapping, drivers,
firmware, Open vSwitch topology, SR-IOV/VF state, and basic connectivity
checks into readable tables and stable JSON.`

// newRootCommand builds the root command and its subtree, binding o to the
// global persistent flags.
func newRootCommand(o *Options) *cobra.Command {
	compat := &compatFlags{}

	root := &cobra.Command{
		Use:           "ifscope",
		Short:         "Unified Linux network interface inspection",
		Long:          longDescription,
		Version:       version.Version,
		SilenceUsage:  true,
		SilenceErrors: true,
		// --netns re-execs the whole process inside the namespace before any
		// command runs, so every view (commands and sysfs) sees that namespace.
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return o.reexecNetns()
		},
		// With no subcommand, legacy view selectors choose the view; otherwise
		// the default is the interfaces+VLANs show view.
		RunE: func(_ *cobra.Command, _ []string) error {
			return o.withWatch(func() error {
				if compat.any() {
					return o.runCompat(*compat)
				}
				return o.runShow()
			})
		},
	}

	o.bindGlobal(root)
	compat.bindCompat(root)
	root.SetVersionTemplate(version.String() + "\n")

	root.AddCommand(
		newShowCommand(o),
		newInterfacesCommand(o),
		newVLANsCommand(o),
		newBondsCommand(o),
		newBridgesCommand(o),
		newBridgeVLANsCommand(o),
		newTunnelsCommand(o),
		newWireGuardCommand(o),
		newPCIeCommand(o),
		newDevlinkCommand(o),
		newRoutesCommand(o),
		newPathCommand(o),
		newRulesCommand(o),
		newNeighborsCommand(o),
		newLLDPCommand(o),
		newFDBCommand(o),
		newDNSCommand(o),
		newOVSCommand(o),
		newSRIOVCommand(o),
		newQdiscCommand(o),
		newClassesCommand(o),
		newFiltersCommand(o),
		newOffloadsCommand(o),
		newQueuesCommand(o),
		newIRQCommand(o),
		newPTPCommand(o),
		newMulticastCommand(o),
		newMDBCommand(o),
		newStatsCommand(o),
		newSocketsCommand(o),
		newNetnsCommand(o),
		newTestCommand(o),
		newAllCommand(o),
		newVersionCommand(),
	)

	// Wrap data subcommands so --watch refreshes them on an interval. version,
	// help, and completion are excluded.
	for _, c := range root.Commands() {
		switch c.Name() {
		case "version", "help", "completion":
			continue
		}
		inner := c.RunE
		if inner == nil {
			continue
		}
		c.RunE = func(cmd *cobra.Command, args []string) error {
			return o.withWatch(func() error { return inner(cmd, args) })
		}
	}

	return root
}

// Execute runs the CLI and returns the process exit code.
func Execute() int {
	o := &Options{}
	root := newRootCommand(o)

	err := root.Execute()
	if err == nil {
		return ExitOK
	}

	var coded *ExitCodeError
	if errors.As(err, &coded) {
		fmt.Fprintln(os.Stderr, "ifscope:", coded.Error())
		return coded.Code
	}

	fmt.Fprintln(os.Stderr, "ifscope:", err)
	return ExitError
}
