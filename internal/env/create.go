package env

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"rpy/internal/python"

	"github.com/pterm/pterm"
	"golang.org/x/term"
)

// CreateProjectVenv creates a local .venv directory in the current working directory.
func CreateProjectVenv(requestedPython string) error {
	info, err := ProjectEnvInfo()
	if err != nil {
		return err
	}

	interpreter, err := resolvePython(requestedPython)
	if err != nil {
		return err
	}

	return createEnvAt(info, interpreter)
}

func CreateSharedVenv(name, requestedPython string, useNow bool) error {
	info, err := SharedEnvInfo(name)
	if err != nil {
		return err
	}

	interpreter, err := resolvePython(requestedPython)
	if err != nil {
		return err
	}

	if err := createEnvAt(info, interpreter); err != nil {
		return err
	}

	if useNow {
		return UseSharedEnv(name)
	}

	return nil
}

func createEnvAt(info Info, interpreter python.Interpreter) error {
	if info.Exists {
		return fmt.Errorf("%s already exists", info.Root)
	}
	if _, err := os.Stat(info.Root); err != nil && !os.IsNotExist(err) {
		return err
	}

	pterm.Println(pterm.FgGray.Sprint("create " + info.Name))
	pterm.Println(
		pterm.FgGray.Sprint("  python: ") +
			pterm.FgLightGreen.Sprint(interpreter.Name) +
			pterm.FgGray.Sprint("  "+shortenPath(interpreter.Path)),
	)
	pterm.Println(
		pterm.FgGray.Sprint("  target: ") +
			pterm.FgGray.Sprint(shortenPath(info.Root)),
	)

	stage := (*pterm.SpinnerPrinter)(nil)
	if term.IsTerminal(int(os.Stdout.Fd())) {
		var spinnerErr error
		stage, spinnerErr = pterm.DefaultSpinner.WithRemoveWhenDone(true).Start("Creating virtual environment")
		if spinnerErr != nil {
			stage = nil
		}
	}
	if stage == nil {
		pterm.Println(pterm.FgGray.Sprint("  creating: ") + info.Root)
	}

	cmd := exec.Command(interpreter.Path, "-m", "venv", info.Root)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if stage != nil {
			stage.Fail()
		}
		message := fmt.Sprintf("failed to create virtual environment: %v", err)
		output := strings.TrimSpace(string(out))
		if output != "" {
			return fmt.Errorf("%s\n%s", message, output)
		}
		return fmt.Errorf("%s", message)
	}

	if stage != nil {
		stage.Stop()
	}
	pterm.Println(pterm.FgLightGreen.Sprint("  OK  ") + shortenPath(info.Root))
	return nil
}

func resolvePython(requested string) (python.Interpreter, error) {
	if requested != "" {
		interpreter, err := python.Resolve(requested)
		if err != nil {
			return python.Interpreter{}, fmt.Errorf("could not resolve requested python interpreter %q: %w", requested, err)
		}
		return interpreter, nil
	}

	interpreter, err := python.ResolveDefault()
	if err != nil {
		return python.Interpreter{}, err
	}

	return interpreter, nil
}

func shortenPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return path
	}

	if path == home {
		return "~"
	}

	prefix := home + string(filepath.Separator)
	if strings.HasPrefix(path, prefix) {
		return "~" + string(filepath.Separator) + strings.TrimPrefix(path, prefix)
	}

	return path
}
