package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	results := make([]Result, 0, 3)

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

	return results, nil
}
