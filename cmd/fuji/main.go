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
	"fuji/internal/sema"
)

func parseBuildCommandArgs(args []string) (src string, out string, noOpt bool, debug bool, err error) {
	var output string
	var file string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--no-opt":
			noOpt = true
		case "--debug":
			debug = true
		case "-o":
			if i+1 >= len(args) {
				return "", "", false, false, fmt.Errorf("-o requires a path")
			}
			i++
			output = args[i]
		default:
			if strings.HasPrefix(args[i], "-") {
				return "", "", false, false, fmt.Errorf("unknown flag: %s", args[i])
			}
			if file != "" {
				return "", "", false, false, fmt.Errorf("multiple source files")
			}
			file = args[i]
		}
	}
	if file == "" {
		return "", "", false, false, fmt.Errorf("usage: fuji build [--no-opt] [--debug] <file.fuji> [-o <exe>]")
	}
	if output == "" {
		output = defaultExeName(file)
	}
	return file, output, noOpt, debug, nil
}

func parseRunCommandArgs(args []string) (src string, noOpt bool, err error) {
	var file string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--no-opt":
			noOpt = true
		default:
			if strings.HasPrefix(args[i], "-") {
				return "", false, fmt.Errorf("unknown flag: %s (fuji run accepts --no-opt)", args[i])
			}
			if file != "" {
				return "", false, fmt.Errorf("multiple source files")
			}
			file = args[i]
		}
	}
	if file == "" {
		return "", false, fmt.Errorf("usage: fuji run [--no-opt] <file.fuji>")
	}
	return file, noOpt, nil
}

// version is set by release builds, e.g. -ldflags "-X main.version=1.0.0"
var version = "0.3.0-dev"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	}

	switch args[0] {
	case "run", "native":
		cmd := args[0]
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "usage: fuji %s [--no-opt] <file.fuji>\n", cmd)
			os.Exit(1)
		}
		path, noOpt, err := parseRunCommandArgs(args[1:])
		if err != nil {
			fatal(err.Error())
		}
		opts := nativebuild.BuildOptions{NoOpt: noOpt}
		if err := api.RunWithBuildOptions(path, "", opts); err != nil {
			fatalErr(err)
		}

	case "watch":
		path, noOpt, err := parseWatchCommandArgs(args[1:])
		if err != nil {
			fatal(err.Error())
		}
		opts := nativebuild.BuildOptions{NoOpt: noOpt}
		if err := runWatch(path, opts); err != nil {
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
		src, output, noOpt, debug, err := parseBuildCommandArgs(args[1:])
		if err != nil {
			fatal(err.Error())
		}
		// --debug prioritizes debuggability (symbols + less optimizer reordering).
		opts := nativebuild.BuildOptions{NoOpt: noOpt || debug, Debug: debug}
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
  run     [--no-opt] <file.fuji>   Compile with LLVM and run (same pipeline as build; --no-opt eases flaky Clang IR opt)
  native  [--no-opt] <file.fuji>   Same as run (backward-compatible alias)
  watch   [--no-opt] <file.fuji>   Rebuild + rerun on .fuji file changes under the entry directory
  check   <file.fuji>              Parse, resolve imports, run sema (no LLVM); prints OK if valid
  fmt     [--check] <files...>     Canonical spacing (4 spaces); ./... walks .fuji files (--check exits 1 if diffs)
  disasm  <file.fuji>              Print LLVM IR for the program (after sema + codegen)
  build   [--no-opt] [--debug] <file.fuji> [-o <exe>]   Native executable (needs llc + Clang; --debug emits symbols and implies --no-opt)
  bundle  <file.fuji> [-o <dir>]   Build + tidy folder to share or sell
  wrap    ...args...                   Forward to fujiwrap (C header → .fuji + wrapper.c); see below
  paths                            Machine-readable toolchain paths (CI / scripts)
  doctor                           Human-readable health check (clang, llc, lld, stdlib, install writable)

  help                             This screen
  version                          Version and platform

OPTIONS
  -o <path>   Output executable (build) or output directory (bundle)

MAINTAINER BUILD (not for end users)
  Building Fuji from source is for contributors only. See CONTRIBUTING.md in the repository.

C / C++ LIBRARIES (readable .fuji wrappers)
  Use the fujiwrap binary from Releases (or next to fuji). Generate bindings from headers:

    fuji wrap -name mylib -headers ./include/mylib.h -out ./wrappers/mylib

  That writes readable mylib.fuji plus wrapper.c and docs. To compile your game:

    set FUJI_NATIVE_SOURCES=wrappers\mylib\wrapper.c
    set FUJI_LINKFLAGS=-I.\include -L.\lib -lmylib
    fuji bundle game.fuji -o dist\mygame

  See docs/wrappers.md and docs/distribution.md.

ENVIRONMENT
  FUJI_CLANG / CC     Clang for native builds (optional if you ship llvm/ next to fuji; see below)
  FUJI_LLC            llc for optional llc+link path (defaults to bundled llvm/bin/llc)
  FUJI_USE_LLD        1 to force -fuse-ld=lld; 0 to disable even if bundled ld.lld exists
  FUJI_PATH           Extra @module search dirs (path list, same separator as PATH)
  FUJI_WRAPPERS       Pre-built .fuji libraries (path list; overrides FUJI_PATH)
  FUJI_NATIVE_SOURCES C/C++ sources linked into your app (e.g. wrapper.c)
  FUJI_LINKFLAGS      Extra linker flags (-lraylib, -L..., frameworks, etc.)
  FUJI_RAYLIB_STAGE   Override path to third_party/.../stage (include/ + lib/) for vendored raylib
  FUJI_USE_VENDORED_RAYLIB  If third_party/raylib_static/stage exists (cwd or next to fuji), prepends -I and links libraylib.a; set 0/false to skip
  FUJI_BUNDLE_FILES   Extra files/dirs copied into the bundle (path-list or quoted paths with spaces)
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
  fuji run --no-opt tests\loops.fuji
  fuji build game.fuji -o game%s
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
		fmt.Println("\nInstall fujiwrap from the same GitHub Releases page as fuji, or place fujiwrap next to fuji.")
		fmt.Println("(Legacy names wrapgen / kujiwrap are still discovered if present.)")
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
	return fmt.Errorf("fujiwrap not found (looked next to fuji and on PATH).\nDownload fujiwrap from GitHub Releases (same page as fuji), or add it next to this executable.\n(legacy binary name wrapgen also works.)")
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
	bundle, err := parser.LoadProgram(path)
	if err != nil {
		return err
	}
	_, err = sema.PrepareNativeBundle(bundle)
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
