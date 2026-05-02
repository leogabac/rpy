package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"rpy/internal/env"
)

type Checker struct {
	// small struct to store the values I care her
	HomeDir string
	RpyDir  string
}

type Result struct {
	// a simple struct for reporting
	Name   string
	OK     bool
	Detail string
}

func NewChecker() (*Checker, error) {
	// simply read the user's home directory, and return a Checker with it
	// this is a simple constructor that returns a Checker pointer
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("resolve home directory: %w", err)
	}

	return &Checker{
		HomeDir: homeDir,
		RpyDir:  filepath.Join(homeDir, ".rpy"),
	}, nil
}

func (c *Checker) Check() ([]Result, error) {
	// then here we make an implementation for the Checker struct

	// first make an empty slice with capacity for three objects
	results := make([]Result, 0, 8)

	// append the home directory that we already have
	results = append(results, Result{
		Name:   "home directory",
		OK:     c.HomeDir != "",
		Detail: c.HomeDir,
	})

	// make the config dir and append the results
	if err := os.MkdirAll(c.RpyDir, 0o755); err != nil {
		return nil, fmt.Errorf("ensure %s exists: %w", c.RpyDir, err)
	}

	results = append(results, Result{
		Name:   "~/.rpy",
		OK:     true,
		Detail: c.RpyDir,
	})

	// see if python exists and append the results
	pythonPath, err := exec.LookPath("python")
	if err != nil {
		results = append(results, Result{
			Name:   "python on PATH",
			OK:     false,
			Detail: "python was not found",
		})
		return results, nil
	}

	results = append(results, Result{
		Name:   "python on PATH",
		OK:     true,
		Detail: pythonPath,
	})

	results = append(results, Result{
		Name:   "shell integration",
		OK:     os.Getenv("RPY_SHELL_INIT") == "1",
		Detail: shellInitDetail(),
	})

	results = append(results, checkCurrentEnvSelection())

	for _, tool := range []string{"make", "cc"} {
		results = append(results, checkTool(tool))
	}

	return results, nil
}

func shellInitDetail() string {
	if os.Getenv("RPY_SHELL_INIT") == "1" {
		return "loaded in current shell session"
	}

	return "not detected in current shell session"
}

func checkCurrentEnvSelection() Result {
	info, err := env.CurrentEnvInfo()
	if err != nil {
		return Result{
			Name:   "current environment",
			OK:     false,
			Detail: err.Error(),
		}
	}

	label := fmt.Sprintf("%s:%s", info.Kind, info.Name)
	if info.Kind == "project" {
		label = "local"
	}
	if info.Exists {
		return Result{
			Name:   "current environment",
			OK:     true,
			Detail: fmt.Sprintf("%s -> %s", label, info.Root),
		}
	}

	return Result{
		Name:   "current environment",
		OK:     false,
		Detail: fmt.Sprintf("%s selected but missing at %s", label, info.Root),
	}
}

func checkTool(name string) Result {
	path, err := exec.LookPath(name)
	if err != nil {
		return Result{
			Name:   name + " available",
			OK:     false,
			Detail: name + " was not found",
		}
	}

	return Result{
		Name:   name + " available",
		OK:     true,
		Detail: path,
	}
}
