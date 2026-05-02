package env

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunInProjectEnv(args []string) error {
	info, err := requireCurrentEnv()
	if err != nil {
		return err
	}

	executable, err := resolveEnvExecutable(info, args[0])
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = envWithProjectBin(info)

	return cmd.Run()
}

func RunPipInProjectEnv(args []string) error {
	info, err := requireCurrentEnv()
	if err != nil {
		return err
	}

	cmd := exec.Command(info.Python, append([]string{"-m", "pip"}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = envWithProjectBin(info)

	return cmd.Run()
}

func requireCurrentEnv() (Info, error) {
	info, err := CurrentEnvInfo()
	if err != nil {
		return Info{}, err
	}

	if !info.Exists {
		if info.Kind == "shared" {
			return Info{}, fmt.Errorf("shared environment %q was selected but not found at %s", info.Name, info.Root)
		}
		return Info{}, fmt.Errorf("no environment found at %s; run `rpy env create` first", info.Root)
	}

	if _, err := os.Stat(info.Python); err != nil {
		if os.IsNotExist(err) {
			return Info{}, fmt.Errorf("%s environment exists but %s is missing", info.Kind, info.Python)
		}

		return Info{}, err
	}

	return info, nil
}

func envWithProjectBin(info Info) []string {
	env := os.Environ()
	pathKey := "PATH="
	pathValue := info.BinDir

	for i, entry := range env {
		if !strings.HasPrefix(entry, pathKey) {
			continue
		}

		current := strings.TrimPrefix(entry, pathKey)
		env[i] = pathKey + pathValue + string(os.PathListSeparator) + current
		return env
	}

	return append(env, pathKey+pathValue)
}

func resolveEnvExecutable(info Info, name string) (string, error) {
	if strings.ContainsRune(name, filepath.Separator) {
		return name, nil
	}

	candidate := filepath.Join(info.BinDir, name)
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}

	path, err := exec.LookPath(name)
	if err != nil {
		return "", err
	}

	return path, nil
}
