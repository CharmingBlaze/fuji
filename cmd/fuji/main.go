package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"fuji/api"
	"fuji/internal/codegen"
	"fuji/internal/diagnostic"
	"fuji/internal/fujihome"
	"fuji/internal/nativebuild"
	"fuji/internal/parser"
)

func parseBuildCommandArgs(args []string) (src string, out string, noOpt bool, err error) {
	var output string
	var file string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--no-opt":
			noOpt = true
		case "-o":
			if i+1 >= len(args) {
				return "", "", false, fmt.Errorf("-o requires a path")
			}
			i++
			output = args[i]
		default:
			if strings.HasPrefix(args[i], "-") {
				return "", "", false, fmt.Errorf("unknown flag: %s", args[i])
			}
			if file != "" {
				return "", "", false, fmt.Errorf("multiple source files")
			}
			file = args[i]
		}
	}
	if file == "" {
		return "", "", false, fmt.Errorf("usage: fuji build [--no-opt] <file.fuji> [-o <exe>]")
	}
	if output == "" {
		output = defaultExeName(file)
	}
	return file, output, noOpt, nil
}

// version is set by release builds, e.g. -ldflags "-X main.version=1.0.0"
var version = "0.2.0-dev"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	}

	switch args[0] {
	case "run", "native":
		cmd := args[0]
		requireArg(args, cmd, "<file.fuji>")
		if err := api.Run(args[1], ""); err != nil {
			fatalErr(err)
		}

	case "check":
		requireArg(args, "check", "<file.fuji>")
		if err := checkFile(args[1]); err != nil {
			fatalErr(err)
		}
		fmt.Println("OK")

	case "fmt":
		if err := runFmtCmd(args[1:]); err != nil {
			fatalErr(err)
		}

	case "disasm":
		requireArg(args, "disasm", "<file.fuji>")
		if err := disasmFile(args[1]); err != nil {
			fatalErr(err)
		}

	case "build":
		src, output, noOpt, err := parseBuildCommandArgs(args[1:])
		if err != nil {
			fatal(err.Error())
		}
		opts := nativebuild.BuildOptions{NoOpt: noOpt}
		if err := buildFileOpts(src, output, opts); err != nil {
			fatalErr(err)
		}

	case "bundle":
		requireArg(args, "bundle", "<file.fuji>")
		outputDir := "dist"
		if len(args) >= 4 && args[2] == "-o" {
			outputDir = args[3]
		}
		if err := bundleFile(args[1], outputDir); err != nil {
			fatalErr(err)
		}

	case "wrap":
		if err := runWrapgen(args[1:]); err != nil {
			fatalErr(err)
		}

	case "paths":
		if err := printResolvedPaths(); err != nil {
			fatalErr(err)
		}

	case "doctor":
		if err := runDoctor(); err != nil {
			fatalErr(err)
		}

	case "version", "--version", "-v":
		fmt.Println(versionLine())

	case "help", "--help", "-h":
		printHelp()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printHelp()
		os.Exit(1)
	}
}

func versionLine() string {
	return fmt.Sprintf("fuji %s (%s/%s)", version, runtime.GOOS, runtime.GOARCH)
}

func printHelp() {
	fmt.Printf(`%s

Fuji compiles .fuji programs to native apps. Your players only need the bundle folder
(executable + assets). They do not need Go, Python, or C++ toolchains.

USAGE
  fuji <command> [arguments]

COMMANDS
  run     <file.fuji>              Compile with LLVM and run the native binary (same pipeline as build)
  native  <file.fuji>              Same as run (backward-compatible alias)
  check   <file.fuji>              Parse + load imports only; prints OK if valid
  fmt     [--check] <files...>     Canonical spacing (4 spaces); ./... walks .fuji files (--check exits 1 if diffs)
  disasm  <file.fuji>              Print LLVM IR for the program (after sema + codegen)
  build   [--no-opt] <file.fuji> [-o <exe>]   Native executable (needs llc + Clang; see CONTRIBUTING.md)
  bundle  <file.fuji> [-o <dir>]   Build + tidy folder to share or sell
  wrap    ...args...                   Forward to fujiwrap (C header → .fuji + wrapper.c); see below
  paths                            Machine-readable toolchain paths (CI / scripts)
  doctor                           Human-readable health check (clang, llc, lld, stdlib, install writable)

  help                             This screen
  version                          Version and platform

OPTIONS
  -o <path>   Output executable (build) or output directory (bundle)

BUILD FROM SOURCE (optional)
  From the repo root with GNU Make; on Windows use mingw32-make when the default make is not GNU make:

    make                # runtime/libfuji_runtime.a, bin/fuji, bin/fujiwrap, then go test
    make runtime-lib    # only the static runtime library
    make raylib-lib     # optional: build Raylib into third_party/raylib_static/stage/ (CMake + raylib/)

C / C++ LIBRARIES (readable .fuji wrappers)
  Parse C headers and emit .fuji bindings plus wrapper.c for the Fuji runtime (Value / fuji_*).
  Library authors build fujiwrap once with Go (or use 'make fujiwrap'):

    go build -o fujiwrap%s ./cmd/wrapgen

  Generate bindings from headers:

    fuji wrap -name mylib -headers ./include/mylib.h -out ./wrappers/mylib

  That writes readable mylib.fuji plus wrapper.c and docs. To compile your game:

    set FUJI_NATIVE_SOURCES=wrappers\mylib\wrapper.c
    set FUJI_LINKFLAGS=-I.\include -L.\lib -lmylib
    fuji bundle game.fuji -o dist\mygame

  See WRAPPERS.md and DISTRIBUTION_GUIDE.md.

ENVIRONMENT
  FUJI_CLANG / CC     Clang for native builds (optional if you ship llvm/ next to fuji; see below)
  FUJI_LLC            llc for optional llc+link path (defaults to bundled llvm/bin/llc)
  FUJI_USE_LLD        1 to force -fuse-ld=lld; 0 to disable even if bundled ld.lld exists
  FUJI_PATH           Extra @module search dirs (path list, same separator as PATH)
  FUJI_WRAPPERS       Pre-built .fuji libraries (path list; overrides FUJI_PATH)
  FUJI_NATIVE_SOURCES C/C++ sources linked into your app (e.g. wrapper.c)
  FUJI_LINKFLAGS      Extra linker flags (-lraylib, -L..., frameworks, etc.)
  FUJI_USE_VENDORED_RAYLIB  If third_party/raylib_static/stage exists (after make raylib-lib), fuji prepends -I and links libraylib.a automatically; set 0/false to skip
  FUJI_BUNDLE_FILES   Extra files copied into the bundle (DLLs, assets)
  FUJI_SKIP_TOOLCHAIN_EXTRACT  If set, never unpack the embedded toolchain archive (dev/CI)
  FUJI_DEBUG_IR       If set, writes .FUJI_build/main.ll in addition to piping IR to clang

SINGLE-EXE / EMBEDDED TOOLCHAIN
  Release builds embed a gzip tarball (see internal/fujihome/embeddata/). On first fuji build
  On first fuji build or fujiwrap run, if the archive lists a real toolchain/bin/clang, it is extracted next to the
  executable. Portable lib/clang/*/include headers can live in that tree; clang is invoked with
  matching -isystem flags. Replace embeddata/bundled_toolchain.tar.gz before go build to ship LLVM.

ZERO-SETUP DISTRIBUTION (same folder as fuji.exe / fuji)
  Put optional trees next to the compiler so users need no separate installs:
    llvm/bin/clang(.exe)   OR   toolchain/bin/clang(.exe)   — used for fuji build and fujiwrap
    stdlib/*.fuji         — @ imports (e.g. @array) resolve here automatically
    wrappers/             — optional generated bindings (same @ resolution rules)
  FUJI_CLANG / FUJI_PATH are only needed when you do not ship these folders or an embedded bundle.

EXAMPLES
  fuji run tests\hello.fuji
  fuji build game.fuji -o game.exe
  fuji bundle game.fuji -o dist\mygame
`, versionLine(), exeExt())
}

func exeExt() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func runWrapgen(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Println("fuji wrap — runs fujiwrap (header → .fuji + wrapper.c). Example:")
		fmt.Println("  fuji wrap -name mylib -headers ./include/mylib.h -out ./wrappers/mylib")
		fmt.Printf("\nBuild fujiwrap and place it next to fuji:\n  go build -o fujiwrap%s ./cmd/wrapgen\n", exeExt())
		fmt.Printf("(legacy names wrapgen / kujiwrap are still discovered if present.)\n")
		return nil
	}
	names := []string{
		"fujiwrap", "fujiwrap.exe",
		"wrapgen", "wrapgen.exe",
		"kujiwrap", "kujiwrap.exe",
	}
	if self, err := os.Executable(); err == nil {
		dir := filepath.Dir(self)
		for _, name := range names {
			p := filepath.Join(dir, name)
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				return runPassthrough(p, args)
			}
		}
	}
	for _, name := range names {
		if p, err := exec.LookPath(name); err == nil {
			return runPassthrough(p, args)
		}
	}
	return fmt.Errorf("fujiwrap not found (looked next to fuji and on PATH).\nBuild once: go build -o fujiwrap%s ./cmd/wrapgen\n(legacy binary name wrapgen also works.)", exeExt())
}

func runPassthrough(path string, args []string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if x, ok := err.(*exec.ExitError); ok && x.ExitCode() != 0 {
			os.Exit(x.ExitCode())
		}
		return err
	}
	return nil
}

func requireArg(args []string, cmd, arg string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: fuji %s %s\n", cmd, arg)
		os.Exit(1)
	}
}

func defaultExeName(source string) string {
	name := strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func checkFile(path string) error {
	_, err := parser.LoadProgram(path)
	return err
}
func disasmFile(path string) error {
	bundle, err := parser.LoadProgram(path)
	if err != nil {
		return err
	}
	ctx, err := codegen.PrepareNativeBundle(bundle)
	if err != nil {
		return err
	}
	mod, err := codegen.EmitLLVMIR(ctx)
	if err != nil {
		return err
	}
	fmt.Print(mod.String())
	return nil
}
func buildFileOpts(path string, output string, opts nativebuild.BuildOptions) error {
	bundle, err := parser.LoadProgram(path)
	if err != nil {
		return err
	}
	return nativebuild.BuildWithOptions(bundle, output, filepath.Base(path), func(s string) { fmt.Print(s) }, opts)
}

func bundleFile(path string, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}
	appName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	exeName := appName
	if runtime.GOOS == "windows" {
		exeName += ".exe"
	}

	fmt.Printf("Bundling %s -> %s/\n\n", filepath.Base(path), outputDir)

	exePath := filepath.Join(outputDir, exeName)
	if err := buildFileOpts(path, exePath, nativebuild.BuildOptions{}); err != nil {
		return err
	}
	if err := writeBundleLauncher(outputDir, exeName); err != nil {
		return err
	}
	if err := writeBundleReadme(outputDir, appName, exeName); err != nil {
		return err
	}
	if err := writeBundleInfo(outputDir, path); err != nil {
		return err
	}
	if err := copyBundleExtraFiles(outputDir); err != nil {
		return err
	}

	fmt.Println()
	fmt.Printf("  Bundle ready: %s\n", outputDir)
	fmt.Printf("  Share the whole folder - friends only run %s (or the launcher script).\n", exeName)
	fmt.Printf("  No Go, Python, or C++ needed on their machine.\n")
	return nil
}

func writeBundleLauncher(outputDir, exeName string) error {
	if runtime.GOOS == "windows" {
		content := fmt.Sprintf("@echo off\r\n\"%%~dp0%s\" %%*\r\n", exeName)
		return os.WriteFile(filepath.Join(outputDir, "run.bat"), []byte(content), 0644)
	}
	content := fmt.Sprintf("#!/bin/sh\nDIR=\"$(CDPATH= cd -- \"$(dirname -- \"$0\")\" && pwd)\"\n\"$DIR/%s\" \"$@\"\n", exeName)
	return os.WriteFile(filepath.Join(outputDir, "run.sh"), []byte(content), 0755)
}

func writeBundleReadme(outputDir, appName, exeName string) error {
	launcher := "run.bat"
	if runtime.GOOS != "windows" {
		launcher = "run.sh"
	}
	lines := []string{
		"# " + appName,
		"",
		"Thanks for trying this app. It was built with **Fuji** - you do not need to install Go, Python, or a C++ compiler.",
		"",
		"## Run",
		"",
		"- **Windows:** double-click `" + launcher + "` or run `" + exeName + "`.",
		"- **macOS / Linux:** in a terminal: `./" + launcher + "` or `./" + exeName + "`",
		"",
		"Keep any DLLs, `.dylib`s, or `assets` folder **in the same folder** as the executable when you move or zip this directory.",
		"",
		"## Files",
		"",
		"| File | What it is |",
		"|------|------------|",
		"| `" + exeName + "` | Your application. |",
		"| `" + launcher + "` | Optional double-click launcher. |",
		"| `README.md` | This file. |",
		"| `bundle-info.txt` | Build notes for developers. |",
		"",
	}
	return os.WriteFile(filepath.Join(outputDir, "README.md"), []byte(strings.Join(lines, "\n")), 0644)
}

func writeBundleInfo(outputDir, sourcePath string) error {
	lines := []string{
		"# Fuji bundle metadata",
		"source=" + sourcePath,
		"native_sources=" + os.Getenv("FUJI_NATIVE_SOURCES"),
		"linkflags=" + os.Getenv("FUJI_LINKFLAGS"),
		"FUJI_version=" + version,
	}
	return os.WriteFile(filepath.Join(outputDir, "bundle-info.txt"), []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func copyBundleExtraFiles(outputDir string) error {
	for _, item := range strings.Fields(os.Getenv("FUJI_BUNDLE_FILES")) {
		dst := filepath.Join(outputDir, filepath.Base(item))
		if err := copyFile(item, dst); err != nil {
			return fmt.Errorf("copy bundle file %s: %w", item, err)
		}
		fmt.Printf("  copied: %s\n", dst)
	}
	return nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// printResolvedPaths prints install-relative toolchain and stdlib resolution for
// debugging portable releases (no FUJI_HOME — layout is always relative to the binary).
func printResolvedPaths() error {
	install, err := fujihome.InstallDir()
	if err != nil {
		return err
	}
	fmt.Printf("FUJI_INSTALL_DIR=%s\n", install)

	clang := fujihome.Clang()
	fmt.Printf("FUJI_CLANG_EFFECTIVE=%s\n", clang)
	fmt.Printf("FUJI_CLANG_STATUS=%s\n", toolResolveStatus(clang))

	llc := fujihome.LLC()
	fmt.Printf("FUJI_LLC_EFFECTIVE=%s\n", llc)
	fmt.Printf("FUJI_LLC_STATUS=%s\n", toolResolveStatus(llc))

	if fujihome.HasBundledLLD() {
		fmt.Println("FUJI_BUNDLED_LLD=1")
	} else {
		fmt.Println("FUJI_BUNDLED_LLD=0")
	}

	stdlib, err := fujihome.StdlibDir()
	if err != nil {
		return err
	}
	fmt.Printf("stdlib_dir=%s\n", stdlib)
	if fi, err := os.Stat(stdlib); err == nil && fi.IsDir() {
		fmt.Println("stdlib_exists=1")
	} else {
		fmt.Println("stdlib_exists=0")
	}

	wrap, err := fujihome.WrappersDir()
	if err != nil {
		return err
	}
	fmt.Printf("wrappers_dir=%s\n", wrap)
	if fi, err := os.Stat(wrap); err == nil && fi.IsDir() {
		fmt.Println("wrappers_exists=1")
	} else {
		fmt.Println("wrappers_exists=0")
	}

	flags, err := fujihome.BundledClangResourceFlags()
	if err != nil {
		fmt.Printf("bundled_clang_include_flags_error=%v\n", err)
	} else {
		fmt.Printf("bundled_clang_include_flags_count=%d\n", len(flags))
	}

	for _, k := range []string{
		"FUJI_CLANG", "FUJI_LLC", "CC", "FUJI_PATH", "FUJI_WRAPPERS",
		"FUJI_SKIP_TOOLCHAIN_EXTRACT", "FUJI_USE_LLD",
	} {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			fmt.Printf("%s=%s\n", k, v)
		}
	}
	fmt.Println("hint: for a human-readable report (writable install dir, tool probes), run: fuji doctor")
	return nil
}

func toolResolveStatus(tool string) string {
	if tool == "" {
		return "empty"
	}
	if filepath.IsAbs(tool) || strings.ContainsRune(tool, filepath.Separator) ||
		(runtime.GOOS == "windows" && strings.ContainsRune(tool, '/')) {
		fi, err := os.Stat(tool)
		if err == nil && !fi.IsDir() {
			return "ok"
		}
		return "missing"
	}
	if _, err := exec.LookPath(tool); err == nil {
		return "on_path"
	}
	return "missing"
}

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "\nerror: %s\n\nRun 'fuji help' for usage.\n", msg)
	os.Exit(1)
}

func fatalErr(err error) {
	fatal(diagnostic.FormatError(err))
}
