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
	Short: "Print metadata about the current environment",
	Long:  `Report whether the currently selected environment exists and show the key paths rpy will use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return env.PrintCurrentEnvInfo()
	},
}

func init() {
	envCmd.AddCommand(infoCmd)
}
