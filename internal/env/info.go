package env

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pterm/pterm"
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
	name, forceLocal, err := currentEnvSelection()
	if err != nil {
		return Info{}, err
	}
	if name != "" && !forceLocal {
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

	pterm.Println(pterm.FgGray.Sprint("env"))
	printInfoLine("selected", displayEnvName(info))
	printInfoLine("root", shortenPath(info.Root))
	printInfoLine("status", envStatus(info.Exists))
	printInfoLine("python", shortenPath(info.Python))
	printInfoLine("pip", shortenPath(info.Pip))
	printInfoLine("activate", shortenPath(info.Activate))

	return nil
}

func printInfoLine(label, value string) {
	pterm.Println(
		pterm.FgGray.Sprint("  "+label+": ") +
			value,
	)
}

func displayEnvName(info Info) string {
	if info.Kind == "shared" {
		return pterm.FgLightGreen.Sprint("shared:" + info.Name)
	}

	return pterm.FgLightGreen.Sprint("local")
}

func envStatus(exists bool) string {
	if exists {
		return pterm.FgLightGreen.Sprint("ready")
	}

	return pterm.FgGray.Sprint("missing")
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
