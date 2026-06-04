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
		// With no subcommand, legacy view selectors choose the view; otherwise
		// the default is the interfaces+VLANs show view.
		RunE: func(_ *cobra.Command, _ []string) error {
			if compat.any() {
				return o.runCompat(*compat)
			}
			return o.runShow()
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
		newRulesCommand(o),
		newNeighborsCommand(o),
		newLLDPCommand(o),
		newFDBCommand(o),
		newDNSCommand(o),
		newOVSCommand(o),
		newSRIOVCommand(o),
		newQdiscCommand(o),
		newClassesCommand(o),
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
