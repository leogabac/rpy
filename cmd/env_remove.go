package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

var envRemoveCmd = &cobra.Command{
	Use:   "remove <name|local>",
	Short: "Remove a shared environment or the local .venv",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return env.RemoveSharedEnv(args[0])
	},
}

func init() {
	envCmd.AddCommand(envRemoveCmd)
}
