package cmd

import "github.com/spf13/cobra"

var pyCmd = &cobra.Command{
	Use:   "py",
	Short: "Manage Python runtimes",
	Long:  `Install, list, and resolve Python runtimes managed by rpy.`,
}

func init() {
	rootCmd.AddCommand(pyCmd)
}
