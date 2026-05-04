package env

import (
	"fmt"
	"os"
)

func PrintUseAndActivateCommand(selection, shell string) error {
	command, err := UseAndActivateCommand(selection, shell)
	if err != nil {
		return err
	}

	fmt.Println(command)
	return nil
}

func UseAndActivateCommand(selection, shell string) (string, error) {
	info, shellValue, err := selectionActivationInfo(selection)
	if err != nil {
		return "", err
	}

	script, err := activationScriptPathForSelection(info, shell)
	if err != nil {
		return "", err
	}

	root := info.Root
	normalized := normalizeShell(shell)
	switch normalized {
	case "fish":
		return fmt.Sprintf(
			"set -gx %s %s; source %q; set -gx RPY_AUTO_VENV %q",
			shellSelectionEnv, shellValue, script, root,
		), nil
	default:
		return fmt.Sprintf(
			"export %s=%q\nsource %q\nexport RPY_AUTO_VENV=%q",
			shellSelectionEnv, shellValue, script, root,
		), nil
	}
}

func selectionActivationInfo(selection string) (Info, string, error) {
	if selection == "local" {
		info, err := ProjectEnvInfo()
		if err != nil {
			return Info{}, "", err
		}
		if !info.Exists {
			return Info{}, "", fmt.Errorf("local environment does not exist at %s", info.Root)
		}
		return info, "local", nil
	}

	info, err := SharedEnvInfo(selection)
	if err != nil {
		return Info{}, "", err
	}
	if !info.Exists {
		return Info{}, "", fmt.Errorf("shared environment %q does not exist at %s", selection, info.Root)
	}
	return info, info.Name, nil
}

func activationScriptPathForSelection(info Info, shell string) (string, error) {
	path := activationScriptPath(info, normalizeShell(shell))
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("activation script not found at %s", path)
		}
		return "", err
	}
	return path, nil
}
