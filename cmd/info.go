/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print metadata about the project-local .venv",
	Long:  `Report whether the .venv exists and show the key paths rpy will use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return env.PrintProjectEnvInfo()
	},
}

func init() {
	envCmd.AddCommand(infoCmd)
}
