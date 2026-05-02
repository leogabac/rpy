package env

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Info struct {
	Name     string
	Kind     string
	Root     string
	BinDir   string
	Python   string
	Pip      string
	Activate string
	Exists   bool
}

func ProjectEnvInfo() (Info, error) {
	root, err := ProjectVenvPath()
	if err != nil {
		return Info{}, err
	}

	info := Info{
		Name:   ".venv",
		Kind:   "project",
		Root:   root,
		BinDir: scriptsDir(root),
	}
	info.Python = filepath.Join(info.BinDir, pythonExecutableName())
	info.Pip = filepath.Join(info.BinDir, pipExecutableName())
	info.Activate = defaultActivationScriptPath(info)

	stat, err := os.Stat(root)
	if err == nil {
		info.Exists = stat.IsDir()
		return info, nil
	}
	if os.IsNotExist(err) {
		return info, nil
	}

	return Info{}, fmt.Errorf("stat %s: %w", root, err)
}

func CurrentEnvInfo() (Info, error) {
	name, err := currentSharedEnvName()
	if err != nil {
		return Info{}, err
	}
	if name != "" {
		return SharedEnvInfo(name)
	}

	return ProjectEnvInfo()
}

func SharedEnvInfo(name string) (Info, error) {
	root, err := SharedEnvPath(name)
	if err != nil {
		return Info{}, err
	}

	info := Info{
		Name:   name,
		Kind:   "shared",
		Root:   root,
		BinDir: scriptsDir(root),
	}
	info.Python = filepath.Join(info.BinDir, pythonExecutableName())
	info.Pip = filepath.Join(info.BinDir, pipExecutableName())
	info.Activate = defaultActivationScriptPath(info)

	stat, err := os.Stat(root)
	if err == nil {
		info.Exists = stat.IsDir()
		return info, nil
	}
	if os.IsNotExist(err) {
		return info, nil
	}

	return Info{}, fmt.Errorf("stat %s: %w", root, err)
}

func PrintCurrentEnvInfo() error {
	info, err := CurrentEnvInfo()
	if err != nil {
		return err
	}

	fmt.Printf("kind: %s\n", info.Kind)
	fmt.Printf("name: %s\n", info.Name)
	fmt.Printf("root: %s\n", info.Root)
	fmt.Printf("exists: %t\n", info.Exists)
	fmt.Printf("bin: %s\n", info.BinDir)
	fmt.Printf("python: %s\n", info.Python)
	fmt.Printf("pip: %s\n", info.Pip)
	fmt.Printf("activate: %s\n", info.Activate)

	return nil
}

func scriptsDir(root string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(root, "Scripts")
	}

	return filepath.Join(root, "bin")
}

func pythonExecutableName() string {
	if runtime.GOOS == "windows" {
		return "python.exe"
	}

	return "python"
}

func pipExecutableName() string {
	if runtime.GOOS == "windows" {
		return "pip.exe"
	}

	return "pip"
}

func defaultActivationScriptPath(info Info) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(info.BinDir, "activate")
	}

	return filepath.Join(info.BinDir, "activate")
}
