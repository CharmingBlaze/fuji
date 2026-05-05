package nativebuild

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"fuji/internal/fujihome"
)

func objFileName() string {
	if runtime.GOOS == "windows" {
		return "main.obj"
	}
	return "main.o"
}

func runLLC(llcPath, irPath, objPath string, optFlag string) error {
	cmd := exec.Command(llcPath, optFlag, "-filetype=obj", "-o", objPath, irPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func linkWithLLDGNU(tc *fujihome.Toolchain, objPath, outAbs string) error {
	args := []string{
		"-flavor", "gnu",
		objPath,
		tc.RuntimeLib,
		"-o", outAbs,
		"-lc", "-lm", "-lpthread",
	}
	cmd := exec.Command(tc.LLD, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func linkWithLLDDarwin(tc *fujihome.Toolchain, objPath, outAbs string) error {
	args := []string{
		"-flavor", "darwin",
		objPath,
		tc.RuntimeLib,
		"-lSystem",
		"-lm",
		"-o", outAbs,
	}
	cmd := exec.Command(tc.LLD, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// runCompileAndLink compiles LLVM IR and links the Fuji runtime (and system libs) in one clang invocation.
func runCompileAndLink(tc *fujihome.Toolchain, irFile, outAbs, rootDir string, opts BuildOptions, log func(string)) error {
	cc := tc.Clang
	if strings.TrimSpace(cc) == "" {
		cc = ClangDriver()
	}
	runtimeInclude := filepath.Join(rootDir, "runtime", "src")

	linkArgs := []string{opts.llcOptFlag()}
	if tc.LLD != "" {
		if runtime.GOOS == "windows" {
			// Clang treats "-fuse-ld=C:\..." as multiple tokens (drive colon); rely on LLVM lld on PATH.
			linkArgs = append(linkArgs, "-fuse-ld=lld")
		} else {
			linkArgs = append(linkArgs, "-fuse-ld="+tc.LLD)
		}
	} else if UseLLD() {
		linkArgs = append(linkArgs, "-fuse-ld=lld")
	}
	if res, err := fujihome.BundledClangResourceFlags(); err == nil {
		linkArgs = append(linkArgs, res...)
	}
	linkArgs = append(linkArgs, irFile, "-I", runtimeInclude)
	nativeSrc := os.Getenv("FUJI_NATIVE_SOURCES")
	if strings.TrimSpace(nativeSrc) == "" {
		nativeSrc = os.Getenv("FUJI_NATIVE_SOURCES")
	}
	if nativeSources := strings.Fields(nativeSrc); len(nativeSources) > 0 {
		if log != nil {
			log(fmt.Sprintf("  native sources: %s\n", strings.Join(nativeSources, " ")))
		}
		linkArgs = append(linkArgs, nativeSources...)
	}
	linkArgs = append(linkArgs, tc.RuntimeLib)
	if inc, arch, ok := vendoredRaylibStatic(rootDir); ok {
		if log != nil {
			log(fmt.Sprintf("  vendored raylib: %s\n", arch))
		}
		linkArgs = append(linkArgs, "-I", inc, arch)
	}
	linkExtra := os.Getenv("FUJI_LINKFLAGS")
	if strings.TrimSpace(linkExtra) == "" {
		linkExtra = os.Getenv("FUJI_LINKFLAGS")
	}
	if extra := strings.Fields(linkExtra); len(extra) > 0 {
		if log != nil {
			log(fmt.Sprintf("  link flags: %s\n\n", strings.Join(extra, " ")))
		}
		linkArgs = append(linkArgs, extra...)
	}
	linkArgs = append(linkArgs, defaultSystemLinkFlags()...)
	if runtime.GOOS == "windows" {
		linkArgs = append(linkArgs, "-lmsvcrt")
	}
	linkArgs = append(linkArgs, "-o", outAbs)

	cmd := exec.Command(cc, linkArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
