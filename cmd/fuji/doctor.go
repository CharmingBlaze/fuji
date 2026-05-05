package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"fuji/internal/fujihome"
)

// runDoctor prints a diagnostic report for zero-install / native builds.
func runDoctor() error {
	fmt.Println("=== fuji doctor ===")
	fmt.Println()

	install, err := fujihome.InstallDir()
	if err != nil {
		return fmt.Errorf("install dir: %w", err)
	}
	fmt.Printf("install_dir: %s\n", install)

	if ok, why := fujihome.InstallDirWritable(install); ok {
		fmt.Println("install_writable: ok")
	} else {
		fmt.Printf("install_writable: NO (%s)\n", why)
		fmt.Println("  hint: move the binary to a normal writable folder; Gatekeeper translocation can block writes.")
	}
	for _, w := range fujihome.InstallDirWarnings(install) {
		fmt.Printf("install_warning: %s\n", w)
	}
	fmt.Println()

	clangPath, clangSrc := fujihome.ClangWithSource()
	fmt.Printf("clang: %s (from %s)\n", clangPath, clangSrc)
	printToolProbe("clang", clangPath)

	llcPath, llcSrc := fujihome.LLCWithSource()
	fmt.Printf("llc: %s (from %s)\n", llcPath, llcSrc)
	printToolProbe("llc", llcPath)

	if p, ok := fujihome.BundledLLDPath(); ok {
		fmt.Printf("lld: %s (bundled)\n", p)
		if fi, err := os.Stat(p); err != nil || fi.IsDir() {
			fmt.Println("lld_status: missing_or_invalid")
		} else {
			fmt.Println("lld_status: ok")
		}
	} else {
		fmt.Println("lld: (not bundled next to fuji; system linker may be used)")
	}
	fmt.Println()

	stdlib, err := fujihome.StdlibDir()
	if err != nil {
		return err
	}
	fmt.Printf("stdlib_dir: %s\n", stdlib)
	printDirStatus("stdlib", stdlib)

	wrap, err := fujihome.WrappersDir()
	if err != nil {
		return err
	}
	fmt.Printf("wrappers_dir: %s\n", wrap)
	printDirStatus("wrappers", wrap)

	flags, err := fujihome.BundledClangResourceFlags()
	if err != nil {
		fmt.Printf("bundled_clang_isystem: error: %v\n", err)
	} else {
		fmt.Printf("bundled_clang_isystem_entries: %d\n", len(flags))
		if len(flags) > 0 && os.Getenv("FUJI_DOCTOR_VERBOSE") != "" {
			fmt.Println(strings.Join(flags, "\n"))
		}
	}

	fmt.Println()
	fmt.Println("Doctor finished. Native builds need clang + llc; LLD is optional if bundled.")
	return nil
}

func printDirStatus(label, path string) {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Printf("%s_exists: 0 (%v)\n", label, err)
		return
	}
	if !fi.IsDir() {
		fmt.Printf("%s_exists: 0 (not a directory)\n", label)
		return
	}
	fmt.Printf("%s_exists: 1\n", label)
}

func printToolProbe(name, path string) {
	if path == "" {
		fmt.Printf("%s_probe: skip (empty path)\n", name)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, "--version")
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s_probe: failed (%v)\n", name, err)
		return
	}
	line := strings.TrimSpace(string(out))
	if idx := strings.IndexByte(line, '\n'); idx >= 0 {
		line = line[:idx]
	}
	if len(line) > 120 {
		line = line[:120] + "..."
	}
	fmt.Printf("%s_probe: %s\n", name, line)
}
