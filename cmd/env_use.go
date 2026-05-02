package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

var envUseCmd = &cobra.Command{
	Use:   "use <name|local>",
	Short: "Select the environment for the current project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] == "local" {
			return env.UseLocalEnv()
		}
		return env.UseSharedEnv(args[0])
	},
}

func init() {
	envCmd.AddCommand(envUseCmd)
}
