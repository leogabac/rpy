package python

import (
	"os"
	"path/filepath"
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
