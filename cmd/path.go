/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

// pathCmd represents the path command
var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the current environment path",
	Long:  `Print the absolute path to the environment currently selected for this project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return env.ResolveEnvPath()
	},
}

func init() {
	envCmd.AddCommand(pathCmd)
}
