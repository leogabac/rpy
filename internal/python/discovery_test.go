package python

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"testing"
)

func TestManagedPythonPathPrefersSourceBuildLayout(t *testing.T) {
	root := t.TempDir()
	binDir := filepath.Join(root, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}

	python3 := filepath.Join(binDir, "python3")
	if err := os.WriteFile(python3, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := managedPythonPath(root)
	if err != nil {
		t.Fatal(err)
	}
	if got != python3 {
		t.Fatalf("managedPythonPath() = %q, want %q", got, python3)
	}
}

func TestManagedPythonPathFallsBackToVersionedBinary(t *testing.T) {
	root := t.TempDir()
	binDir := filepath.Join(root, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}

	python312 := filepath.Join(binDir, "python3.12")
	if err := os.WriteFile(python312, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	got, err := managedPythonPath(root)
	if err != nil {
		t.Fatal(err)
	}
	if got != python312 {
		t.Fatalf("managedPythonPath() = %q, want %q", got, python312)
	}
}

func TestDiscoverIncludesManagedAndPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	managedRoot := filepath.Join(home, ".rpy", "pythons", "3.12.0", "bin")
	if err := os.MkdirAll(managedRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	managedPython := filepath.Join(managedRoot, "python3")
	if err := os.WriteFile(managedPython, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	pathDir := filepath.Join(t.TempDir(), "bin")
	if err := os.MkdirAll(pathDir, 0o755); err != nil {
		t.Fatal(err)
	}
	pathPython := filepath.Join(pathDir, "python3.11")
	if err := os.WriteFile(pathPython, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", pathDir)

	items, err := Discover()
	if err != nil {
		t.Fatal(err)
	}

	paths := make([]string, 0, len(items))
	for _, item := range items {
		paths = append(paths, item.Path)
	}

	if !slices.Contains(paths, managedPython) {
		t.Fatalf("Discover() missing managed python %q", managedPython)
	}
	if !slices.Contains(paths, pathPython) {
		t.Fatalf("Discover() missing PATH python %q", pathPython)
	}
}

func TestResolveDefaultFallsBackToManaged(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("PATH", filepath.Join(t.TempDir(), "empty"))

	managedRoot := filepath.Join(home, ".rpy", "pythons", "3.12.0", "bin")
	if err := os.MkdirAll(managedRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	managedPython := filepath.Join(managedRoot, "python3")
	if err := os.WriteFile(managedPython, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	item, err := ResolveDefault()
	if err != nil {
		t.Fatal(err)
	}
	if item.Path != managedPython {
		t.Fatalf("ResolveDefault() = %q, want %q", item.Path, managedPython)
	}
}

func TestCompareVersionish(t *testing.T) {
	if got := compareVersionish("3.12.0", "3.11.9"); got <= 0 {
		t.Fatalf("compareVersionish() = %d, want > 0", got)
	}
	if got := compareVersionish("python3.12", "python3.9"); got <= 0 {
		t.Fatalf("compareVersionish() = %d, want > 0", got)
	}
	if runtime.GOOS == "" {
		t.Fatal("unreachable")
	}
}

func TestDisplayInterpreterNamePrefersResolvedVersionedName(t *testing.T) {
	got := displayInterpreterName("python", "/usr/bin/python3.14")
	if got != "python3.14" {
		t.Fatalf("displayInterpreterName() = %q, want %q", got, "python3.14")
	}

	got = displayInterpreterName("python3", "/usr/bin/python3.12")
	if got != "python3.12" {
		t.Fatalf("displayInterpreterName() = %q, want %q", got, "python3.12")
	}

	got = displayInterpreterName("python3.11", "/usr/bin/python3.11")
	if got != "python3.11" {
		t.Fatalf("displayInterpreterName() = %q, want %q", got, "python3.11")
	}
}
