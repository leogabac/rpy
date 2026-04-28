/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import "github.com/spf13/cobra"

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:                "install <packages...>",
	Short:              "Install packages into the project-local environment",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cobra.MinimumNArgs(1)(cmd, args)
		}
		return runPipWithSubcommand("install", args)
	},
}

func init() {
	pipCmd.AddCommand(installCmd)
}
