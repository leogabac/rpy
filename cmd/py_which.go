package cmd

import (
	"rpy/internal/python"

	"github.com/spf13/cobra"
)

var pyWhichCmd = &cobra.Command{
	Use:   "which <name-or-version>",
	Short: "Resolve a managed Python runtime or explicit interpreter path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return python.PrintResolved(args[0])
	},
}

func init() {
	pyCmd.AddCommand(pyWhichCmd)
}
