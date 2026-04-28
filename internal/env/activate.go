package env

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func PrintActivateCommand(shell string) error {
	info, err := requireProjectEnv()
	if err != nil {
		return err
	}

	command, err := ActivateCommand(info, shell)
	if err != nil {
		return err
	}

	fmt.Println(command)
	return nil
}

func PrintActivationPath(shell string) error {
	info, err := requireProjectEnv()
	if err != nil {
		return err
	}

	fmt.Println(activationScriptPath(info, normalizeShell(shell)))
	return nil
}

func ActivateCommand(info Info, shell string) (string, error) {
	normalized := normalizeShell(shell)
	script := activationScriptPath(info, normalized)

	if _, err := os.Stat(script); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("activation script not found at %s", script)
		}
		return "", err
	}

	switch normalized {
	case "fish":
		return fmt.Sprintf("source %q", script), nil
	case "powershell", "pwsh":
		return fmt.Sprintf("& %q", script), nil
	case "cmd":
		return script, nil
	default:
		return fmt.Sprintf("source %q", script), nil
	}
}

func activationScriptPath(info Info, shell string) string {
	if runtime.GOOS == "windows" {
		switch shell {
		case "powershell", "pwsh":
			return filepath.Join(info.BinDir, "Activate.ps1")
		case "cmd":
			return filepath.Join(info.BinDir, "activate.bat")
		default:
			return filepath.Join(info.BinDir, "activate")
		}
	}

	switch shell {
	case "fish":
		return filepath.Join(info.BinDir, "activate.fish")
	case "csh":
		return filepath.Join(info.BinDir, "activate.csh")
	default:
		return filepath.Join(info.BinDir, "activate")
	}
}

func normalizeShell(shell string) string {
	if shell == "" {
		shell = filepath.Base(os.Getenv("SHELL"))
	}

	shell = strings.ToLower(shell)
	switch shell {
	case "", "bash", "sh":
		return "sh"
	case "zsh":
		return "zsh"
	case "fish":
		return "fish"
	case "powershell", "pwsh":
		return "powershell"
	case "cmd", "cmd.exe":
		return "cmd"
	case "csh", "tcsh":
		return "csh"
	default:
		return shell
	}
}
