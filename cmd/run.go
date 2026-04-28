/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:                "run <command> [args...]",
	Short:              "Run a command inside the project-local .venv",
	Long:               `Run a command with the environment's bin directory prepended to PATH.`,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cobra.MinimumNArgs(1)(cmd, args)
		}
		return env.RunInProjectEnv(args)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
