package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the local and shared environments available to this project",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return env.PrintEnvs()
	},
}

func init() {
	envCmd.AddCommand(envListCmd)
}
