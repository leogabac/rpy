package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

var activateShell string
var activatePathOnly bool

var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Print the shell command needed to activate the current environment",
	Long: `Print the activation snippet for the current environment.

Because a child process cannot modify the parent shell environment, this
command prints the shell command you should evaluate, for example:

  eval "$(rpy env activate)"
  source "$(rpy env activate --path)"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if activatePathOnly {
			return env.PrintActivationPath(activateShell)
		}

		return env.PrintActivateCommand(activateShell)
	},
}

func init() {
	envCmd.AddCommand(activateCmd)
	activateCmd.Flags().StringVar(&activateShell, "shell", "", "Shell to target: bash, zsh, fish, powershell, pwsh, cmd")
	activateCmd.Flags().BoolVar(&activatePathOnly, "path", false, "Print only the activation script path")
	_ = activateCmd.Flags().MarkHidden("path")
}
