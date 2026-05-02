package cmd

import (
	"rpy/internal/python"

	"github.com/spf13/cobra"
)

var pyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List managed and PATH Python runtimes",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return python.PrintDiscovered()
	},
}

func init() {
	pyCmd.AddCommand(pyListCmd)
}
