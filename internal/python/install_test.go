package python

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSourceArchiveURL(t *testing.T) {
	got := sourceArchiveURL("3.12.3")
	want := "https://www.python.org/ftp/python/3.12.3/Python-3.12.3.tgz"
	if got != want {
		t.Fatalf("sourceArchiveURL() = %q, want %q", got, want)
	}
}

func TestSafeJoinRejectsEscape(t *testing.T) {
	if _, err := safeJoin("/tmp/root", "../etc/passwd"); err == nil {
		t.Fatal("expected path escape error")
	}
}

func TestExtractSourceTarGz(t *testing.T) {
	tmp := t.TempDir()
	archive := filepath.Join(tmp, "Python-3.12.3.tgz")
	if err := writeTestArchive(archive); err != nil {
		t.Fatal(err)
	}

	dest := filepath.Join(tmp, "out")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatal(err)
	}

	root, err := extractSourceTarGz(archive, dest)
	if err != nil {
		t.Fatal(err)
	}

	if root != filepath.Join(dest, "Python-3.12.3") {
		t.Fatalf("extractSourceTarGz() root = %q", root)
	}
	if _, err := os.Stat(filepath.Join(root, "README")); err != nil {
		t.Fatalf("expected extracted source file: %v", err)
	}
}

func TestDefaultMakeJobs(t *testing.T) {
	if got := defaultMakeJobs(8); got != 8 {
		t.Fatalf("defaultMakeJobs(8) = %d", got)
	}
	if got := defaultMakeJobs(0); got < 1 || got > runtime.NumCPU() {
		t.Fatalf("defaultMakeJobs(0) = %d", got)
	}
}

func TestReadLogTail(t *testing.T) {
	path := filepath.Join(t.TempDir(), "build.log")
	if err := os.WriteFile(path, []byte("a\nb\nc\nd\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got := readLogTail(path, 2)
	if got != "c\nd" {
		t.Fatalf("readLogTail() = %q", got)
	}
}

func writeTestArchive(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gzw := gzip.NewWriter(f)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	if err := tw.WriteHeader(&tar.Header{
		Name:     "Python-3.12.3/",
		Typeflag: tar.TypeDir,
		Mode:     0o755,
	}); err != nil {
		return err
	}

	contents := []byte("cpython source tree")
	if err := tw.WriteHeader(&tar.Header{
		Name:     "Python-3.12.3/README",
		Typeflag: tar.TypeReg,
		Mode:     0o644,
		Size:     int64(len(contents)),
	}); err != nil {
		return err
	}

	_, err = tw.Write(contents)
	return err
}
