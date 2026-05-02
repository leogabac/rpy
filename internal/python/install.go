package python

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

const pythonSourceBaseURL = "https://www.python.org/ftp/python"

type InstallOptions struct {
	FromFile string
	Force    bool
	Jobs     int
}

type managedMetadata struct {
	Version     string    `json:"version"`
	Method      string    `json:"method"`
	SourceURL   string    `json:"source_url,omitempty"`
	ArchivePath string    `json:"archive_path,omitempty"`
	InstalledAt time.Time `json:"installed_at"`
}

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

type installUI struct{}

func Install(version string, opts InstallOptions) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("managed installs currently support Linux only")
	}

	ui := installUI{}
	ui.header(version)

	destRoot, err := managedInstallDir(version)
	if err != nil {
		return err
	}

	if err := prepareDestination(destRoot, opts.Force); err != nil {
		return err
	}

	archivePath := opts.FromFile
	sourceURL := ""
	if archivePath == "" {
		sourceURL = sourceArchiveURL(version)
		stage := ui.startStage("Fetching source archive")
		stage.updateDetail(sourceURL)
		archivePath, err = downloadSource(version, sourceURL)
		if err != nil {
			stage.fail("Fetch failed")
			return err
		}
		stage.success("Source archive ready")
	} else {
		ui.info("Using local source archive", archivePath)
	}

	return installFromSource(ui, version, archivePath, destRoot, opts.Jobs, managedMetadata{
		Version:     version,
		Method:      "source-build",
		SourceURL:   sourceURL,
		ArchivePath: archivePath,
		InstalledAt: time.Now().UTC(),
	})
}

func managedRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".rpy", "pythons"), nil
}

func downloadsRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".rpy", "downloads"), nil
}

func logsRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".rpy", "logs"), nil
}

func managedInstallDir(version string) (string, error) {
	root, err := managedRoot()
	if err != nil {
		return "", err
	}

	return filepath.Join(root, version), nil
}

func prepareDestination(destRoot string, force bool) error {
	if _, err := os.Stat(destRoot); err == nil {
		if !force {
			return fmt.Errorf("%s already exists; use --force to replace it", destRoot)
		}
		if err := os.RemoveAll(destRoot); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(destRoot), 0o755); err != nil {
		return err
	}

	return nil
}

func sourceArchiveURL(version string) string {
	return fmt.Sprintf("%s/%s/Python-%s.tgz", pythonSourceBaseURL, version, version)
}

func downloadSource(version, url string) (string, error) {
	root, err := downloadsRoot()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return "", err
	}

	dest := filepath.Join(root, "Python-"+version+".tgz")
	if _, err := os.Stat(dest); err == nil {
		return dest, nil
	} else if !os.IsNotExist(err) {
		return "", err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "rpy")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download %s: %s", url, resp.Status)
	}

	tmp := dest + ".part"
	f, err := os.Create(tmp)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		_ = os.Remove(tmp)
		return "", err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return "", err
	}

	if err := os.Rename(tmp, dest); err != nil {
		_ = os.Remove(tmp)
		return "", err
	}

	return dest, nil
}

func installFromSource(ui installUI, version, archivePath, destRoot string, jobs int, metadata managedMetadata) error {
	stageRoot := destRoot + ".tmp"
	_ = os.RemoveAll(stageRoot)
	if err := os.MkdirAll(stageRoot, 0o755); err != nil {
		return err
	}

	logDir, err := createInstallLogDir(version)
	if err != nil {
		return err
	}

	success := false
	defer func() {
		_ = os.RemoveAll(stageRoot)
		if success {
			_ = os.RemoveAll(logDir)
		}
	}()

	stage := ui.startStage("Extracting source archive")
	stage.updateDetail(archivePath)
	srcRoot, err := extractSourceTarGz(archivePath, stageRoot)
	if err != nil {
		stage.fail("Extraction failed")
		return err
	}
	stage.success("Source extracted")

	stage = ui.startStage("Configuring build")
	stage.updateDetail(srcRoot)
	if err := runBuildStep(srcRoot, filepath.Join(logDir, "configure.log"), "./configure", "--prefix="+destRoot, "--with-ensurepip=install"); err != nil {
		stage.fail("Configure failed")
		return wrapBuildError("configure", err, filepath.Join(logDir, "configure.log"))
	}
	stage.success("Configure complete")

	makeJobs := defaultMakeJobs(jobs)
	stage = ui.startStage("Compiling Python")
	stage.updateDetail(fmt.Sprintf("%d parallel job(s)", makeJobs))
	if err := runBuildStep(srcRoot, filepath.Join(logDir, "make.log"), "make", "-j", strconv.Itoa(makeJobs)); err != nil {
		stage.fail("Compilation failed")
		return wrapBuildError("make", err, filepath.Join(logDir, "make.log"))
	}
	stage.success("Compilation complete")

	stage = ui.startStage("Installing runtime")
	stage.updateDetail(destRoot)
	if err := runBuildStep(srcRoot, filepath.Join(logDir, "make-install.log"), "make", "install"); err != nil {
		stage.fail("Install failed")
		return wrapBuildError("make install", err, filepath.Join(logDir, "make-install.log"))
	}
	stage.success("Runtime installed")

	if _, err := os.Stat(filepath.Join(destRoot, "bin", "python3")); err != nil {
		return fmt.Errorf("install completed without %s: %w", filepath.Join(destRoot, "bin", "python3"), err)
	}

	if err := writeMetadata(destRoot, metadata); err != nil {
		return err
	}

	success = true
	ui.success("Managed interpreter ready", filepath.Join(destRoot, "bin", "python3"))
	return nil
}

func defaultMakeJobs(explicit int) int {
	if explicit > 0 {
		return explicit
	}
	if runtime.NumCPU() < 1 {
		return 1
	}
	return runtime.NumCPU()
}

func createInstallLogDir(version string) (string, error) {
	root, err := logsRoot()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return "", err
	}

	dir := filepath.Join(root, "install-"+version+"-"+time.Now().UTC().Format("20060102T150405Z"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	return dir, nil
}

func runBuildStep(dir, logPath string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	logFile, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer logFile.Close()

	cmd.Stdout = logFile
	cmd.Stderr = logFile
	return cmd.Run()
}

func wrapBuildError(step string, err error, logPath string) error {
	tail := readLogTail(logPath, 20)
	if tail != "" {
		return fmt.Errorf("%s failed: %w\nlog: %s\nrecent output:\n%s\ninstall the required system build dependencies yourself first; rpy does not manage them", step, err, logPath, tail)
	}

	return fmt.Errorf("%s failed: %w\nlog: %s\ninstall the required system build dependencies yourself first; rpy does not manage them", step, err, logPath)
}

func readLogTail(path string, lines int) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	chunks := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(chunks) > lines {
		chunks = chunks[len(chunks)-lines:]
	}
	return strings.Join(chunks, "\n")
}

func (installUI) header(version string) {
	pterm.DefaultSection.Println("rpy py install " + version)
}

func (installUI) info(label, detail string) {
	pterm.Println(styleMuted("  "+label+": ") + detail)
}

func (installUI) success(label, detail string) {
	pterm.Println(styleOK("  OK  ") + label + ": " + detail)
}

func (installUI) startStage(title string) installStage {
	spinner, err := pterm.DefaultSpinner.WithRemoveWhenDone(true).Start(title)
	if err != nil {
		pterm.Printf("  -> %s\n", title)
		return installStage{title: title}
	}

	return installStage{title: title, spinner: spinner}
}

type installStage struct {
	title   string
	spinner *pterm.SpinnerPrinter
}

func (s installStage) updateDetail(detail string) {
	if strings.TrimSpace(detail) == "" {
		return
	}
	if s.spinner != nil {
		s.spinner.UpdateText(fmt.Sprintf("%s: %s", s.title, detail))
		return
	}
	pterm.Println(styleMuted("     "+s.title+": ") + detail)
}

func (s installStage) success(message string) {
	if s.spinner != nil {
		s.spinner.Stop()
		pterm.Println(styleOK("  OK  ") + message)
		return
	}
	pterm.Println(styleOK("  OK  ") + message)
}

func (s installStage) fail(message string) {
	if s.spinner != nil {
		s.spinner.Fail()
		pterm.Println(styleFail("  !!  ") + message)
		return
	}
	pterm.Println(styleFail("  !!  ") + message)
}

func styleMuted(text string) string {
	return pterm.FgGray.Sprint(text)
}

func styleOK(text string) string {
	return pterm.FgLightGreen.Sprint(text)
}

func styleFail(text string) string {
	return pterm.FgRed.Sprint(text)
}

func extractSourceTarGz(archivePath, dest string) (string, error) {
	if err := extractTarGz(archivePath, dest); err != nil {
		return "", err
	}

	entries, err := os.ReadDir(dest)
	if err != nil {
		return "", err
	}
	if len(entries) != 1 || !entries[0].IsDir() {
		return "", fmt.Errorf("expected a single source directory inside %s", archivePath)
	}

	return filepath.Join(dest, entries[0].Name()), nil
}

func extractTarGz(archivePath, dest string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		targetPath, err := safeJoin(dest, hdr.Name)
		if err != nil {
			return err
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(hdr.Mode)); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return err
			}
			out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			if err := out.Close(); err != nil {
				return err
			}
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
				return err
			}
			if err := os.Symlink(hdr.Linkname, targetPath); err != nil {
				return err
			}
		default:
			// Python source archives are simple enough that unsupported types can be skipped for now.
		}
	}
}

func safeJoin(root, name string) (string, error) {
	cleanName := filepath.Clean(name)
	target := filepath.Join(root, cleanName)
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("archive entry %q escaped destination", name)
	}
	return target, nil
}

func writeMetadata(dest string, metadata managedMetadata) error {
	payload, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dest, "rpy.json"), payload, 0o644)
}
