/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rpy",
	Short: "Manage Python runtimes and local or shared virtual environments",
	Long: `rpy is a small Python workflow wrapper for Linux-first development.

It manages:
- Python runtimes under ~/.rpy/pythons
- a local project environment at ./.venv
- shared named environments under ~/.rpy/envs

Each project uses either its local environment or one selected shared environment.`,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}

func init() {}
