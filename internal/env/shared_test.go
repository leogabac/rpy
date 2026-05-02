package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCurrentEnvInfoUsesSharedSelection(t *testing.T) {
	home := t.TempDir()
	project := t.TempDir()
	t.Setenv("HOME", home)

	restore := chdirForTest(t, project)
	defer restore()

	sharedPython := filepath.Join(home, ".rpy", "envs", "glass", "bin", "python")
	if err := os.MkdirAll(filepath.Dir(sharedPython), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sharedPython, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(project, projectSelectionFile), []byte("glass\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	info, err := CurrentEnvInfo()
	if err != nil {
		t.Fatal(err)
	}
	if info.Kind != "shared" || info.Name != "glass" {
		t.Fatalf("CurrentEnvInfo() kind/name = %q/%q", info.Kind, info.Name)
	}
	if info.Root != filepath.Join(home, ".rpy", "envs", "glass") {
		t.Fatalf("CurrentEnvInfo() root = %q", info.Root)
	}
}

func TestCurrentEnvInfoFallsBackToProjectEnv(t *testing.T) {
	project := t.TempDir()

	restore := chdirForTest(t, project)
	defer restore()

	info, err := CurrentEnvInfo()
	if err != nil {
		t.Fatal(err)
	}
	if info.Kind != "project" || info.Name != ".venv" {
		t.Fatalf("CurrentEnvInfo() kind/name = %q/%q", info.Kind, info.Name)
	}
	if info.Root != filepath.Join(project, ".venv") {
		t.Fatalf("CurrentEnvInfo() root = %q", info.Root)
	}
}

func TestUseSharedEnvWritesProjectSelection(t *testing.T) {
	home := t.TempDir()
	project := t.TempDir()
	t.Setenv("HOME", home)

	restore := chdirForTest(t, project)
	defer restore()

	sharedPython := filepath.Join(home, ".rpy", "envs", "glass", "bin", "python")
	if err := os.MkdirAll(filepath.Dir(sharedPython), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sharedPython, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := UseSharedEnv("glass"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(project, projectSelectionFile))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "glass\n" {
		t.Fatalf("selection file = %q", string(data))
	}
}

func TestUseLocalEnvRemovesProjectSelection(t *testing.T) {
	project := t.TempDir()

	restore := chdirForTest(t, project)
	defer restore()

	if err := os.WriteFile(filepath.Join(project, projectSelectionFile), []byte("glass\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := UseLocalEnv(); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(project, projectSelectionFile)); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be removed, got %v", projectSelectionFile, err)
	}
}

func chdirForTest(t *testing.T, dir string) func() {
	t.Helper()

	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	return func() {
		if err := os.Chdir(old); err != nil {
			t.Fatal(err)
		}
	}
}
