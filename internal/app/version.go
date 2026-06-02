package app

import (
	"fmt"

	"github.com/bgrewell/ifscope/internal/version"
	"github.com/spf13/cobra"
)

// newVersionCommand prints detailed build metadata.
func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), version.String())
			return nil
		},
	}
}
