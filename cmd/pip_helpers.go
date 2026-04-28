package cmd

import "rpy/internal/env"

func runPipWithSubcommand(name string, args []string) error {
	fullArgs := append([]string{name}, args...)
	return env.RunPipInProjectEnv(fullArgs)
}
