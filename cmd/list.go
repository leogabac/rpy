package cmd

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:                "list [pip-list-args...]",
	Short:              "List packages installed in the current environment",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPipWithSubcommand("list", args)
	},
}

func init() {
	pipCmd.AddCommand(listCmd)
}
