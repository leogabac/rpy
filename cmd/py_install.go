package cmd

import (
	"rpy/internal/python"

	"github.com/spf13/cobra"
)

var pyInstallFromFile string
var pyInstallForce bool
var pyInstallJobs int

var pyInstallCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Build and install a managed Python runtime from source",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return python.Install(args[0], python.InstallOptions{
			FromFile: pyInstallFromFile,
			Force:    pyInstallForce,
			Jobs:     pyInstallJobs,
		})
	},
}

func init() {
	pyCmd.AddCommand(pyInstallCmd)
	pyInstallCmd.Flags().StringVar(&pyInstallFromFile, "from-file", "", "Build from a local Python source tarball instead of downloading it")
	pyInstallCmd.Flags().BoolVar(&pyInstallForce, "force", false, "Replace an existing managed runtime with the same version")
	pyInstallCmd.Flags().IntVar(&pyInstallJobs, "jobs", 0, "Number of parallel make jobs to use; defaults to CPU count")
}
