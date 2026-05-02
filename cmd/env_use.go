package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

var envUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Select a shared environment for the current project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return env.UseSharedEnv(args[0])
	},
}

func init() {
	envCmd.AddCommand(envUseCmd)
}
