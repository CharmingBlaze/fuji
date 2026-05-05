package fujihome

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// LinkMode selects how the native backend links the object file with the runtime.
type LinkMode int

const (
	// LinkClang uses the Clang driver (same resolution as [Clang]) with the usual flags.
	LinkClang LinkMode = iota
	// LinkLLDGNU uses LLVM's LLD in GNU flavor (Linux embedded releases).
	LinkLLDGNU
	// LinkLLDDarwin uses LLVM's LLD in Darwin / ld64-compatible flavor (macOS embedded releases).
	LinkLLDDarwin
)

// Toolchain holds absolute paths to LLVM tools and the Fuji static archive used by [nativebuild].
type Toolchain struct {
	LLC        string
	LLD        string // used when LinkMode == LinkLLDGNU; may be empty for LinkClang
	Clang      string // used when LinkMode == LinkClang
	RuntimeLib string
	LinkMode   LinkMode
}

// FindToolchain resolves LLVM + runtime for native builds.
// When built with "-tags release" and a populated internal/fujihome/bundled tree, the
// embedded layout is extracted once to a temp directory and returned.
// Otherwise the toolchain comes from the environment and PATH (see [ClangWithSource], [LLCWithSource]).
func FindToolchain() (*Toolchain, error) {
	tc, err := embeddedToolchain()
	if err != nil {
		return nil, err
	}
	if tc != nil {
		return tc, nil
	}
	return findSystemToolchain()
}

func findSystemToolchain() (*Toolchain, error) {
	root, err := filepath.Abs(".")
	if err != nil {
		return nil, fmt.Errorf("project root: %w", err)
	}
	llc, _ := LLCWithSource()
	clang := Clang()
	lld := findLLDBinary()
	return &Toolchain{
		LLC:        llc,
		LLD:        lld,
		Clang:      clang,
		RuntimeLib: filepath.Join(root, "runtime", "libfuji_runtime.a"),
		LinkMode:   LinkClang,
	}, nil
}

func findLLDBinary() string {
	candidates := []string{"ld.lld-14", "ld.lld-15", "ld.lld"}
	if runtime.GOOS == "windows" {
		candidates = []string{"ld.lld.exe", "lld.exe"}
	}
	for _, name := range candidates {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}
	return ""
}

// ErrIncompleteEmbeddedToolchain means a release build was compiled with -tags release
// but bundled/{GOOS}/{GOARCH} is missing one of llc, lld, or libfuji_runtime.a.
var ErrIncompleteEmbeddedToolchain = errors.New("incomplete embedded toolchain (see internal/fujihome/bundled/README.md)")

func findClangForEmbeddedLink(llcPath string) string {
	if dir := filepath.Dir(llcPath); dir != "" {
		name := "clang"
		if runtime.GOOS == "windows" {
			name = "clang.exe"
		}
		candidate := filepath.Join(dir, name)
		if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
			return candidate
		}
	}
	return Clang()
}

func platformKey() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}

func llcBinaryName() string {
	if runtime.GOOS == "windows" {
		return "llc.exe"
	}
	return "llc"
}

func lldBinaryName() string {
	if runtime.GOOS == "windows" {
		return "lld.exe"
	}
	return "lld"
}
