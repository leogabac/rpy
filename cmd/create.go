/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

var createPython string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a project-local .venv",
	Long:  `Create a .venv in the current working directory using Python's built-in venv module.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return env.CreateProjectVenv(createPython)
	},
}

func init() {
	envCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createPython, "python", "", "Python interpreter to use")
}
