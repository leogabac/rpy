/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

// pipCmd represents the pip command
var pipCmd = &cobra.Command{
	Use:                "pip <pip-args...>",
	Short:              "Run pip through the project-local environment interpreter",
	Long:               `Invoke python -m pip using the interpreter from the current project's .venv.`,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cobra.MinimumNArgs(1)(cmd, args)
		}
		return env.RunPipInProjectEnv(args)
	},
}

func init() {
	rootCmd.AddCommand(pipCmd)
}
