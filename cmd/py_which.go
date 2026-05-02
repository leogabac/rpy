package cmd

import (
	"rpy/internal/python"

	"github.com/spf13/cobra"
)

var pyWhichCmd = &cobra.Command{
	Use:   "which <name-or-version>",
	Short: "Resolve a Python runtime from managed installs, PATH, or an explicit path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return python.PrintResolved(args[0])
	},
}

func init() {
	pyCmd.AddCommand(pyWhichCmd)
}
