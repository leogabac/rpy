package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List local and shared environments",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return env.PrintEnvs()
	},
}

func init() {
	envCmd.AddCommand(envListCmd)
}
