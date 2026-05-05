package api

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"fuji/internal/nativebuild"
)

// Run compiles the entry .fuji with the native (LLVM) pipeline, writes a temporary
// executable, and runs it. Optional overlay replaces on-disk source for the entry
// file only. Stdout/stderr of the build steps and the child process default to os.Stdout/os.Stderr.
func Run(path, overlay string) error {
	return RunWithWriters(path, overlay, os.Stdout, os.Stderr)
}

// RunWithWriters is like Run but sends compiler log output and the executed program's
// stdout/stderr to the given writers (e.g. for embedding or IDEs).
func RunWithWriters(path, overlay string, stdout, stderr io.Writer) error {
	absEntry, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	overlays := map[string]string{}
	if overlay != "" {
		overlays[absEntry] = overlay
	}
	tmp, err := tempExecutablePath()
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tmp) }()

	log := func(s string) {
		_, _ = io.WriteString(stdout, s)
	}
	if err := nativebuild.BuildWithOverlays(path, overlays, tmp, log); err != nil {
		return err
	}
	cmd := exec.Command(tmp)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func tempExecutablePath() (string, error) {
	pattern := "fuji_run_*"
	if runtime.GOOS == "windows" {
		pattern = "fuji_run_*.exe"
	}
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	name := f.Name()
	_ = f.Close()
	if err := os.Remove(name); err != nil {
		return "", err
	}
	return name, nil
}
