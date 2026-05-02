/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

var createPython string
var createUse bool

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a project-local or shared virtual environment",
	Long:  `Create a .venv in the current working directory, or create a shared named environment when a name is provided.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 0:
			return env.CreateProjectVenv(createPython)
		case 1:
			return env.CreateSharedVenv(args[0], createPython, createUse)
		default:
			return cobra.MaximumNArgs(1)(cmd, args)
		}
	},
}

func init() {
	envCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createPython, "python", "", "Python interpreter to use")
	createCmd.Flags().BoolVar(&createUse, "use", false, "When creating a shared environment, select it for the current project")
	createCmd.Flags().SetInterspersed(true)
}
