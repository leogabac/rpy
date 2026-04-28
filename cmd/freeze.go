package cmd

import "github.com/spf13/cobra"

var freezeCmd = &cobra.Command{
	Use:                "freeze [pip-freeze-args...]",
	Short:              "Print pinned packages from the project-local environment",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPipWithSubcommand("freeze", args)
	},
}

func init() {
	pipCmd.AddCommand(freezeCmd)
}
