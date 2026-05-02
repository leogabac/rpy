package cmd

import (
	"rpy/internal/python"

	"github.com/spf13/cobra"
)

var pyRemoveCmd = &cobra.Command{
	Use:   "remove <version>",
	Short: "Remove a managed Python runtime",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return python.Remove(args[0])
	},
}

func init() {
	pyCmd.AddCommand(pyRemoveCmd)
}
