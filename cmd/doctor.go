package cmd

import (
	"fmt"

	"rpy/internal/doctor"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check whether the local rpy setup looks healthy",
	Long: `doctor runs a small set of environment checks and prints a short report.

The first version should stay simple:
- confirm the home directory can be resolved
- confirm ~/.rpy can be created
- confirm python is available on PATH
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		checker, err := doctor.NewChecker()
		if err != nil {
			return err
		}

		report, err := checker.Check()
		if err != nil {
			return err
		}

		for _, item := range report {
			status := "ok"
			if !item.OK {
				status = "fail"
			}

			if item.Detail != "" {
				fmt.Printf("[%s] %s: %s\n", status, item.Name, item.Detail)
				continue
			}

			fmt.Printf("[%s] %s\n", status, item.Name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
