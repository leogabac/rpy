package python

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/pterm/pterm"
)

type Interpreter struct {
	Name   string
	Path   string
	Source string
}

var pythonNamePattern = regexp.MustCompile(`^python([0-9]+(\.[0-9]+)*)?$`)

func PrintDiscovered() error {
	interpreters, err := Discover()
	if err != nil {
		return err
	}

	if len(interpreters) == 0 {
		pterm.Warning.Println("No Python interpreters found")
		return nil
	}

	printInterpreterGroup("Managed", filterInterpreters(interpreters, "managed"))
	printInterpreterGroup("PATH", filterInterpreters(interpreters, "PATH"))

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
	managed, err := discoverManaged()
	if err != nil {
		return nil, err
	}
	pathItems := discoverPath()

	items := append(make([]Interpreter, 0, len(managed)+len(pathItems)), managed...)
	items = append(items, pathItems...)

	slices.SortFunc(items, compareInterpreters)
	return dedupeInterpreters(items), nil
}

func Resolve(query string) (Interpreter, error) {
	items, err := Discover()
	if err != nil {
		return Interpreter{}, err
	}

	if strings.ContainsRune(query, filepath.Separator) {
		if _, err := os.Stat(query); err != nil {
			return Interpreter{}, err
		}
		return Interpreter{Name: filepath.Base(query), Path: query, Source: "path"}, nil
	}

	if path, err := exec.LookPath(query); err == nil {
		canonical := normalizeInterpreterPath(path)
		return Interpreter{Name: displayInterpreterName(filepath.Base(path), canonical), Path: canonical, Source: "PATH"}, nil
	}

	for _, item := range items {
		if matches(item, query) {
			return item, nil
		}
	}

	return Interpreter{}, fmt.Errorf("no python interpreter matched %q", query)
}

func ResolveDefault() (Interpreter, error) {
	for _, name := range []string{"python", "python3"} {
		if path, err := exec.LookPath(name); err == nil {
			canonical := normalizeInterpreterPath(path)
			return Interpreter{
				Name:   displayInterpreterName(filepath.Base(path), canonical),
				Path:   canonical,
				Source: "PATH",
			}, nil
		}
	}

	items, err := Discover()
	if err != nil {
		return Interpreter{}, err
	}
	if len(items) == 0 {
		return Interpreter{}, fmt.Errorf("could not find a usable python interpreter on PATH or under ~/.rpy/pythons")
	}

	return items[0], nil
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

func discoverPath() []Interpreter {
	pathValue := strings.TrimSpace(os.Getenv("PATH"))
	if pathValue == "" {
		return nil
	}

	seen := make(map[string]bool)
	items := make([]Interpreter, 0, 8)
	for _, dir := range filepath.SplitList(pathValue) {
		if dir == "" {
			continue
		}

		matches, err := filepath.Glob(filepath.Join(dir, "python*"))
		if err != nil {
			continue
		}

		for _, candidate := range matches {
			base := filepath.Base(candidate)
			if !pythonNamePattern.MatchString(base) {
				continue
			}

			info, err := os.Stat(candidate)
			if err != nil || info.IsDir() {
				continue
			}
			if info.Mode()&0o111 == 0 {
				continue
			}
			canonical := normalizeInterpreterPath(candidate)
			if seen[canonical] {
				continue
			}

			seen[canonical] = true
			items = append(items, Interpreter{
				Name:   displayInterpreterName(base, canonical),
				Path:   canonical,
				Source: "PATH",
			})
		}
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

func compareInterpreters(a, b Interpreter) int {
	if a.Source != b.Source {
		if a.Source == "managed" {
			return -1
		}
		if b.Source == "managed" {
			return 1
		}
	}

	if diff := compareVersionish(a.Name, b.Name); diff != 0 {
		return -diff
	}

	if a.Name != b.Name {
		return strings.Compare(a.Name, b.Name)
	}

	return strings.Compare(a.Path, b.Path)
}

func compareVersionish(a, b string) int {
	aParts, aOK := numericParts(a)
	bParts, bOK := numericParts(b)
	if !aOK || !bOK {
		return 0
	}

	limit := len(aParts)
	if len(bParts) > limit {
		limit = len(bParts)
	}

	for i := 0; i < limit; i++ {
		var av, bv int
		if i < len(aParts) {
			av = aParts[i]
		}
		if i < len(bParts) {
			bv = bParts[i]
		}
		if av < bv {
			return -1
		}
		if av > bv {
			return 1
		}
	}

	return 0
}

func numericParts(name string) ([]int, bool) {
	trimmed := strings.TrimPrefix(name, "python")
	trimmed = strings.Trim(trimmed, ".")
	if trimmed == "" {
		return nil, false
	}

	rawParts := strings.Split(trimmed, ".")
	parts := make([]int, 0, len(rawParts))
	for _, part := range rawParts {
		if part == "" {
			return nil, false
		}

		value := 0
		for _, ch := range part {
			if ch < '0' || ch > '9' {
				return nil, false
			}
			value = value*10 + int(ch-'0')
		}
		parts = append(parts, value)
	}

	return parts, true
}

func dedupeInterpreters(items []Interpreter) []Interpreter {
	seen := make(map[string]bool, len(items))
	out := make([]Interpreter, 0, len(items))
	for _, item := range items {
		key := normalizeInterpreterPath(item.Path)
		item.Path = key
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, item)
	}
	return out
}

func normalizeInterpreterPath(path string) string {
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path
	}
	return resolved
}

func displayInterpreterName(name, canonicalPath string) string {
	canonicalBase := filepath.Base(canonicalPath)
	if canonicalBase == "" {
		return name
	}
	if name == "python" || name == "python3" {
		if pythonNamePattern.MatchString(canonicalBase) {
			return canonicalBase
		}
	}
	return name
}

func filterInterpreters(items []Interpreter, source string) []Interpreter {
	out := make([]Interpreter, 0, len(items))
	for _, item := range items {
		if item.Source == source {
			out = append(out, item)
		}
	}
	return out
}

func printInterpreterGroup(title string, items []Interpreter) {
	if len(items) == 0 {
		return
	}

	pterm.Println(pterm.FgGray.Sprint(strings.ToLower(title)))
	for _, item := range items {
		pterm.Println(
			pterm.FgLightGreen.Sprint("  "+item.Name) +
				pterm.FgGray.Sprint("  "+shortenPath(item.Path)),
		)
	}
}

func shortenPath(path string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return path
	}

	if path == home {
		return "~"
	}

	prefix := home + string(filepath.Separator)
	if strings.HasPrefix(path, prefix) {
		return "~" + string(filepath.Separator) + strings.TrimPrefix(path, prefix)
	}

	return path
}
