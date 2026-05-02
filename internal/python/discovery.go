package python

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Interpreter struct {
	Name   string
	Path   string
	Source string
}

func PrintDiscovered() error {
	interpreters, err := Discover()
	if err != nil {
		return err
	}

	for _, item := range interpreters {
		fmt.Printf("%s\t%s\t%s\n", item.Name, item.Source, item.Path)
	}

	return nil
}

func PrintResolved(query string) error {
	item, err := Resolve(query)
	if err != nil {
		return err
	}

	fmt.Println(item.Path)
	return nil
}

func Discover() ([]Interpreter, error) {
	items, err := discoverManaged()
	if err != nil {
		return nil, err
	}
	return items, nil
}

func Resolve(query string) (Interpreter, error) {
	managed, err := discoverManaged()
	if err != nil {
		return Interpreter{}, err
	}
	for _, item := range managed {
		if matches(item, query) {
			return item, nil
		}
	}

	if strings.ContainsRune(query, filepath.Separator) {
		if _, err := os.Stat(query); err != nil {
			return Interpreter{}, err
		}
		return Interpreter{Name: filepath.Base(query), Path: query, Source: "path"}, nil
	}

	if path, err := exec.LookPath(query); err == nil {
		return Interpreter{Name: filepath.Base(path), Path: path, Source: "PATH"}, nil
	}

	return Interpreter{}, fmt.Errorf("no python interpreter matched %q", query)
}

func discoverManaged() ([]Interpreter, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	root := filepath.Join(home, ".rpy", "pythons")
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	items := make([]Interpreter, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		candidate, err := managedPythonPath(filepath.Join(root, entry.Name()))
		if err != nil {
			continue
		}

		items = append(items, Interpreter{
			Name:   entry.Name(),
			Path:   candidate,
			Source: "managed",
		})
	}

	return items, nil
}

func managedPythonPath(root string) (string, error) {
	candidates := []string{
		filepath.Join(root, "bin", "python"),
		filepath.Join(root, "bin", "python3"),
	}

	matches, err := filepath.Glob(filepath.Join(root, "bin", "python3.*"))
	if err != nil {
		return "", err
	}
	candidates = append(candidates, matches...)

	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err != nil || info.IsDir() {
			continue
		}
		return candidate, nil
	}

	return "", os.ErrNotExist
}
func matches(item Interpreter, query string) bool {
	if item.Name == query {
		return true
	}
	if item.Path == query {
		return true
	}
	if strings.TrimPrefix(item.Name, "python") == query {
		return true
	}
	if strings.Contains(item.Path, query) {
		return true
	}
	return false
}
