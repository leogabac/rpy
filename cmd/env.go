/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import "github.com/spf13/cobra"

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage the project-local virtual environment",
	Long:  `Manage the .venv associated with the current project directory.`,
}

func init() {
	rootCmd.AddCommand(envCmd)
}
