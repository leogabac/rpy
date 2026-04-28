package env

import (
	"fmt"
	"os"
	"os/exec"
)

// CreateProjectVenv creates a local .venv directory in the current working directory.
func CreateProjectVenv(python string) error {
	info, err := ProjectEnvInfo()
	if err != nil {
		return err
	}

	if info.Exists {
		return fmt.Errorf("%s already exists", info.Root)
	}
	if _, err := os.Stat(info.Root); err != nil && !os.IsNotExist(err) {
		return err
	}

	pythonBin, err := findPython(python)
	if err != nil {
		return err
	}

	cmd := exec.Command(pythonBin, "-m", "venv", info.Root)

	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		fmt.Print(string(out))
	}
	if err != nil {
		return fmt.Errorf("failed to create virtual environment: %w", err)
	}

	fmt.Println(info.Root)
	return nil
}

func findPython(requested string) (string, error) {
	if requested != "" {
		path, err := exec.LookPath(requested)
		if err != nil {
			return "", fmt.Errorf("could not find requested python interpreter %q on PATH", requested)
		}

		return path, nil
	}

	for _, name := range []string{"python", "python3"} {
		path, err := exec.LookPath(name)
		if err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("could not find python or python3 on PATH")
}
