package cmd

import "github.com/spf13/cobra"

var pyCmd = &cobra.Command{
	Use:   "py",
	Short: "Discover Python interpreters",
	Long:  `List and resolve Python interpreters from managed runtimes and PATH.`,
}

func init() {
	rootCmd.AddCommand(pyCmd)
}
