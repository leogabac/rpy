package env

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProjectRoot returns the effective project root.
// If the current directory is inside a Git repository, the repository root is used.
// Otherwise the current working directory is treated as the project root.
func ProjectRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	root, err := findGitRoot(cwd)
	if err != nil {
		return "", err
	}
	if root != "" {
		return root, nil
	}

	return cwd, nil
}

// ProjectVenvPath returns the path to the project-local virtual environment.
func ProjectVenvPath() (string, error) {
	root, err := ProjectRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(root, ".venv"), nil
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

func findGitRoot(start string) (string, error) {
	current := start
	for {
		gitPath := filepath.Join(current, ".git")
		ok, err := isGitWorktreeMarker(gitPath)
		if err != nil {
			return "", err
		}
		if ok {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", nil
		}
		current = parent
	}
}

func isGitWorktreeMarker(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if info.Mode().IsRegular() {
		return true, nil
	}
	if !info.IsDir() {
		return false, nil
	}

	head := filepath.Join(path, "HEAD")
	if _, err := os.Stat(head); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
