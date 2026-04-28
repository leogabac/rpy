package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

var shellInitShell string
var shellInitInstall bool

var shellInitCmd = &cobra.Command{
	Use:   "shell-init",
	Short: "Print shell integration for activation and auto-activation",
	Long: `Print shell integration code for your shell startup file.

Example:

  eval "$(rpy shell-init --shell bash)"

After loading the integration, rpy env activate can update the current shell
session and rpy can auto-activate a local .venv when you enter a project
directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if shellInitInstall {
			return env.PrintShellInstallHint(shellInitShell)
		}

		return env.PrintShellInit(shellInitShell)
	},
}

func init() {
	rootCmd.AddCommand(shellInitCmd)
	shellInitCmd.Flags().StringVar(&shellInitShell, "shell", "", "Shell to target: bash, zsh, fish")
	shellInitCmd.Flags().BoolVar(&shellInitInstall, "install", false, "Print the startup-file line to install shell integration")
}
