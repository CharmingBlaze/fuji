//go:build FUJI_wip
// +build FUJI_wip

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// fujiBin is the path to a freshly built fuji executable (see TestMain).
var fujiBin string

// moduleRoot is the repository root (directory containing go.mod).
var moduleRoot string

func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	moduleRoot = filepath.Clean(filepath.Join(wd, "..", ".."))

	tmp, err := os.MkdirTemp("", "fuji-cli")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() { _ = os.RemoveAll(tmp) }()

	name := "fuji"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	fujiBin = filepath.Join(tmp, name)

	build := exec.Command("go", "build", "-o", fujiBin, ".")
	build.Dir = wd
	build.Env = os.Environ()
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "go build: %v\n%s", err, out)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func scriptPath(parts ...string) string {
	return filepath.Join(append([]string{moduleRoot}, parts...)...)
}

func TestCLI_noArgsShowsHelp(t *testing.T) {
	cmd := exec.Command(fujiBin)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("no-args: %v\n%s", err, out)
	}
	if cmd.ProcessState.ExitCode() != 0 {
		t.Fatalf("exit code = %d want 0 (print help)", cmd.ProcessState.ExitCode())
	}
	if !bytes.Contains(out, []byte("USAGE")) || !bytes.Contains(out, []byte("COMMANDS")) {
		t.Fatalf("expected help text, got:\n%s", out)
	}
}

func TestCLI_version(t *testing.T) {
	cmd := exec.Command(fujiBin, "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version: %v\n%s", err, out)
	}
	if !bytes.Contains(out, []byte("fuji")) {
		t.Fatalf("expected version line, got %q", string(out))
	}
}

func TestCLI_wrap_help(t *testing.T) {
	cmd := exec.Command(fujiBin, "wrap", "--help")
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("wrap: %v\n%s", err, out)
	}
	if !bytes.Contains(out, []byte("wrapgen")) {
		t.Fatalf("expected wrap help mentioning wrapgen, got:\n%s", out)
	}
}

func TestCLI_check_hello(t *testing.T) {
	cmd := exec.Command(fujiBin, "check", scriptPath("tests", "hello.fuji"))
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("check: %v\n%s", err, out)
	}
	if string(out) != "OK\n" {
		t.Fatalf("stdout = %q want OK + newline", string(out))
	}
}

func TestCLI_check_missingFile(t *testing.T) {
	cmd := exec.Command(fujiBin, "check", scriptPath("tests", "does_not_exist.fuji"))
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if cmd.ProcessState.ExitCode() != 1 {
		t.Fatalf("exit code = %d want 1", cmd.ProcessState.ExitCode())
	}
	if len(out) == 0 {
		t.Fatal("expected error message on stderr/stdout")
	}
}

func TestCLI_run_hello(t *testing.T) {
	cmd := exec.Command(fujiBin, "run", scriptPath("tests", "hello.fuji"))
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	if !bytes.Contains(out, []byte("Hello, Fuji!")) {
		t.Fatalf("output missing greeting:\n%s", out)
	}
}

func TestCLI_run_importModule(t *testing.T) {
	cmd := exec.Command(fujiBin, "run", scriptPath("tests", "import_test.fuji"))
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	if !bytes.Contains(out, []byte("25")) {
		t.Fatalf("output missing square result:\n%s", out)
	}
}

func TestCLI_run_phase1Surface(t *testing.T) {
	cmd := exec.Command(fujiBin, "run", scriptPath("tests", "phase1_surface.fuji"))
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	// Checksum includes for-in over array elements (values), not indices; see internal/fuji/compiler.go desugarForInStmt.
	if !bytes.Contains(out, []byte("98")) {
		t.Fatalf("output missing expected checksum:\n%s", out)
	}
}

func TestCLI_disasm_hello(t *testing.T) {
	cmd := exec.Command(fujiBin, "disasm", scriptPath("tests", "hello.fuji"))
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("disasm: %v\n%s", err, out)
	}
	if !bytes.Contains(out, []byte("OP_")) {
		t.Fatalf("disasm output missing opcodes:\n%s", out)
	}
}

func TestCLI_build_hello(t *testing.T) {
	if _, err := exec.LookPath("clang"); err != nil {
		t.Skip("LLVM clang not on PATH:", err)
	}

	tmpDir, err := os.MkdirTemp("", "fuji-build")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	outName := "hello_out"
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}
	outPath := filepath.Join(tmpDir, outName)

	cmd := exec.Command(fujiBin, "build", scriptPath("tests", "hello.fuji"), "-o", outPath)
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}

	run := exec.Command(outPath)
	run.Dir = tmpDir
	runOut, err := run.CombinedOutput()
	if err != nil {
		t.Fatalf("run built binary: %v\n%s", err, runOut)
	}
	if !bytes.Contains(runOut, []byte("Hello, Fuji!")) {
		t.Fatalf("built program output missing greeting:\n%s", runOut)
	}
}

func normalizeStdoutNewlines(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func TestCLI_run_and_build_print_spacing(t *testing.T) {
	want := "hello world\na b c\n\n"

	runVM := exec.Command(fujiBin, "run", scriptPath("tests", "print_spacing.fuji"))
	runVM.Dir = moduleRoot
	vmOut, err := runVM.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, vmOut)
	}
	if normalizeStdoutNewlines(string(vmOut)) != want {
		t.Fatalf("VM print output mismatch\ngot:  %q\nwant: %q", string(vmOut), want)
	}

	if _, err := exec.LookPath("clang"); err != nil {
		t.Skip("LLVM clang not on PATH:", err)
	}

	tmpDir, err := os.MkdirTemp("", "fuji-print-build")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	outName := "print_out"
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}
	outPath := filepath.Join(tmpDir, outName)

	build := exec.Command(fujiBin, "build", scriptPath("tests", "print_spacing.fuji"), "-o", outPath)
	build.Dir = tmpDir
	bout, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, bout)
	}

	runNative := exec.Command(outPath)
	runNative.Dir = tmpDir
	nativeOut, err := runNative.CombinedOutput()
	if err != nil {
		t.Fatalf("run native: %v\n%s", err, nativeOut)
	}
	if normalizeStdoutNewlines(string(nativeOut)) != want {
		t.Fatalf("native print output mismatch VM\ngot:  %q\nwant: %q", string(nativeOut), want)
	}
}

func TestCLI_run_and_build_default_rest(t *testing.T) {
	want := "3\n3\n"

	runVM := exec.Command(fujiBin, "run", scriptPath("tests", "default_rest.fuji"))
	runVM.Dir = moduleRoot
	vmOut, err := runVM.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, vmOut)
	}
	if normalizeStdoutNewlines(string(vmOut)) != want {
		t.Fatalf("VM output mismatch\ngot:  %q\nwant: %q", string(vmOut), want)
	}

	if _, err := exec.LookPath("clang"); err != nil {
		t.Skip("LLVM clang not on PATH:", err)
	}

	tmpDir, err := os.MkdirTemp("", "fuji-default-rest")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	outName := "defrest_out"
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}
	outPath := filepath.Join(tmpDir, outName)

	build := exec.Command(fujiBin, "build", scriptPath("tests", "default_rest.fuji"), "-o", outPath)
	build.Dir = tmpDir
	bout, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, bout)
	}

	runNative := exec.Command(outPath)
	runNative.Dir = tmpDir
	nativeOut, err := runNative.CombinedOutput()
	if err != nil {
		t.Fatalf("run native: %v\n%s", err, nativeOut)
	}
	if normalizeStdoutNewlines(string(nativeOut)) != want {
		t.Fatalf("native output mismatch VM\ngot:  %q\nwant: %q", string(nativeOut), want)
	}
}

func TestCLI_run_and_build_debug_minimal(t *testing.T) {
	runVM := exec.Command(fujiBin, "run", scriptPath("debug_minimal.fuji"))
	runVM.Dir = moduleRoot
	vmOut, err := runVM.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, vmOut)
	}
	want := normalizeStdoutNewlines(string(vmOut))

	if _, err := exec.LookPath("clang"); err != nil {
		t.Skip("LLVM clang not on PATH:", err)
	}

	tmpDir, err := os.MkdirTemp("", "fuji-debug-minimal")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	outName := "dbgmin_out"
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}
	outPath := filepath.Join(tmpDir, outName)

	build := exec.Command(fujiBin, "build", scriptPath("debug_minimal.fuji"), "-o", outPath)
	build.Dir = tmpDir
	bout, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, bout)
	}

	runNative := exec.Command(outPath)
	runNative.Dir = tmpDir
	nativeOut, err := runNative.CombinedOutput()
	if err != nil {
		t.Fatalf("run native: %v\n%s", err, nativeOut)
	}
	if normalizeStdoutNewlines(string(nativeOut)) != want {
		t.Fatalf("native output mismatch VM\ngot:  %q\nwant: %q", string(nativeOut), want)
	}
}

func TestCLI_run_and_build_control_fuji(t *testing.T) {
	runVM := exec.Command(fujiBin, "run", scriptPath("tests", "control.fuji"))
	runVM.Dir = moduleRoot
	vmOut, err := runVM.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, vmOut)
	}
	want := normalizeStdoutNewlines(string(vmOut))

	if _, err := exec.LookPath("clang"); err != nil {
		t.Skip("LLVM clang not on PATH:", err)
	}

	tmpDir, err := os.MkdirTemp("", "fuji-control-build")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	outName := "control_out"
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}
	outPath := filepath.Join(tmpDir, outName)

	build := exec.Command(fujiBin, "build", scriptPath("tests", "control.fuji"), "-o", outPath)
	build.Dir = tmpDir
	bout, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, bout)
	}

	runNative := exec.Command(outPath)
	runNative.Dir = tmpDir
	nativeOut, err := runNative.CombinedOutput()
	if err != nil {
		t.Fatalf("run native: %v\n%s", err, nativeOut)
	}
	if normalizeStdoutNewlines(string(nativeOut)) != want {
		t.Fatalf("native output mismatch VM\ngot:  %q\nwant: %q", string(nativeOut), want)
	}
}
