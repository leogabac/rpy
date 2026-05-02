package env

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pterm/pterm"
)

const projectSelectionFile = ".rpy-env"

var sharedEnvNamePattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

func SharedEnvPath(name string) (string, error) {
	name, err := validateSharedEnvName(name)
	if err != nil {
		return "", err
	}

	root, err := sharedEnvRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(root, name), nil
}

func ListSharedEnvs() ([]Info, error) {
	root, err := sharedEnvRoot()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	items := make([]Info, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		info, err := SharedEnvInfo(entry.Name())
		if err != nil {
			return nil, err
		}
		items = append(items, info)
	}

	return items, nil
}

func PrintEnvs() error {
	project, err := ProjectEnvInfo()
	if err != nil {
		return err
	}

	items, err := ListSharedEnvs()
	if err != nil {
		return err
	}

	selected, err := currentSharedEnvName()
	if err != nil {
		return err
	}

	currentKind := "project"
	if selected != "" {
		currentKind = "shared"
	}

	pterm.Println(pterm.FgGray.Sprint("local"))
	projectPrefix := "  "
	if currentKind == "project" {
		projectPrefix = "* "
	}
	projectStatus := ""
	if !project.Exists {
		projectStatus = pterm.FgGray.Sprint(" (missing)")
	}
	pterm.Println(
		pterm.FgLightGreen.Sprint(projectPrefix+project.Name) +
			pterm.FgGray.Sprint("  "+shortenPath(project.Root)) +
			projectStatus,
	)

	if len(items) == 0 {
		return nil
	}

	pterm.Println(pterm.FgGray.Sprint("shared"))
	for _, item := range items {
		prefix := "  "
		if item.Name == selected {
			prefix = "* "
		}
		pterm.Println(
			pterm.FgLightGreen.Sprint(prefix+item.Name) +
				pterm.FgGray.Sprint("  "+shortenPath(item.Root)),
		)
	}

	return nil
}

func UseSharedEnv(name string) error {
	info, err := SharedEnvInfo(name)
	if err != nil {
		return err
	}
	if !info.Exists {
		return fmt.Errorf("shared environment %q does not exist at %s", name, info.Root)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(cwd, projectSelectionFile)
	return os.WriteFile(path, []byte(info.Name+"\n"), 0o644)
}

func RemoveSharedEnv(name string) error {
	path, err := SharedEnvPath(name)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("shared environment %q does not exist", name)
		}
		return err
	}

	return os.RemoveAll(path)
}

func currentSharedEnvName() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(filepath.Join(cwd, projectSelectionFile))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	name := strings.TrimSpace(string(data))
	if name == "" {
		return "", nil
	}

	return validateSharedEnvName(name)
}

func sharedEnvRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".rpy", "envs"), nil
}

func validateSharedEnvName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("environment name cannot be empty")
	}
	if !sharedEnvNamePattern.MatchString(name) {
		return "", fmt.Errorf("invalid environment name %q", name)
	}
	return name, nil
}
