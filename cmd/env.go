/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import "github.com/spf13/cobra"

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage local and shared virtual environments",
	Long:  `Manage the project's local .venv and any shared named environments under ~/.rpy/envs.`,
}

func init() {
	rootCmd.AddCommand(envCmd)
}
