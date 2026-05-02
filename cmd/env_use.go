package cmd

import (
	"rpy/internal/env"

	"github.com/spf13/cobra"
)

var envUseActivate bool

var envUseCmd = &cobra.Command{
	Use:   "use <name|local>",
	Short: "Select the environment for the current project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] == "local" {
			if err := env.UseLocalEnv(); err != nil {
				return err
			}
		} else {
			if err := env.UseSharedEnv(args[0]); err != nil {
				return err
			}
		}

		if envUseActivate {
			return env.PrintActivateCommand("")
		}

		return nil
	},
}

func init() {
	envCmd.AddCommand(envUseCmd)
	envUseCmd.Flags().BoolVar(&envUseActivate, "activate", false, "Print the activation command after switching environments")
}
