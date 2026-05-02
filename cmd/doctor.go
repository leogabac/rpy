package cmd

import (
	"rpy/internal/doctor"

	"github.com/pterm/pterm"
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

		pterm.Println(pterm.FgGray.Sprint("doctor"))
		for _, item := range report {
			prefix := pterm.FgLightGreen.Sprint("  OK  ")
			if !item.OK {
				prefix = pterm.FgRed.Sprint("  !!  ")
			}
			if item.Detail != "" {
				pterm.Println(prefix + item.Name + pterm.FgGray.Sprint(": "+item.Detail))
				continue
			}

			pterm.Println(prefix + item.Name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
