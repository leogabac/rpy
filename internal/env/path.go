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

// ResolveEnvPath prints the path to the project-local virtual environment.
func ResolveEnvPath() error {
	envPath, err := ProjectVenvPath()
	if err != nil {
		return err
	}

	fmt.Println(envPath)
	return nil
}
