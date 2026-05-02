package env

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProjectVenvPath returns the path to the project-local virtual environment.
func ProjectVenvPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, ".venv"), nil
}

// ResolveEnvPath prints the path to the currently selected environment.
func ResolveEnvPath() error {
	info, err := CurrentEnvInfo()
	if err != nil {
		return err
	}

	fmt.Println(info.Root)
	return nil
}
