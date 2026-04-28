package python

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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
	seen := make(map[string]bool)
	items := make([]Interpreter, 0, 16)

	managed, err := discoverManaged()
	if err != nil {
		return nil, err
	}
	for _, item := range managed {
		if seen[item.Path] {
			continue
		}
		seen[item.Path] = true
		items = append(items, item)
	}

	pathItems := discoverPath()
	for _, item := range pathItems {
		if seen[item.Path] {
			continue
		}
		seen[item.Path] = true
		items = append(items, item)
	}

	slices.SortFunc(items, func(a, b Interpreter) int {
		return strings.Compare(a.Name+"\t"+a.Path, b.Name+"\t"+b.Path)
	})

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

	for _, item := range discoverPath() {
		if matches(item, query) {
			return item, nil
		}
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

		candidate := filepath.Join(root, entry.Name(), "bin", "python")
		if _, err := os.Stat(candidate); err == nil {
			items = append(items, Interpreter{
				Name:   entry.Name(),
				Path:   candidate,
				Source: "managed",
			})
		}
	}

	return items, nil
}

func discoverPath() []Interpreter {
	candidates := []string{
		"python",
		"python3",
		"python3.9",
		"python3.10",
		"python3.11",
		"python3.12",
		"python3.13",
		"python3.14",
	}

	items := make([]Interpreter, 0, len(candidates))
	for _, name := range candidates {
		path, err := exec.LookPath(name)
		if err != nil {
			continue
		}

		items = append(items, Interpreter{
			Name:   name,
			Path:   path,
			Source: "PATH",
		})
	}

	return items
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
