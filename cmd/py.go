package cmd

import "github.com/spf13/cobra"

var pyCmd = &cobra.Command{
	Use:   "py",
	Short: "Manage Python runtimes",
	Long:  `Install, list, and resolve Python runtimes from rpy-managed installs and PATH.`,
}

func init() {
	rootCmd.AddCommand(pyCmd)
}
