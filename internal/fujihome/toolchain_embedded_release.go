//go:build release

package fujihome

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

//go:embed all:bundled
var bundledReleaseFS embed.FS

var (
	embeddedOnce sync.Once
	embeddedTC   *Toolchain
	embeddedErr  error
)

func embeddedToolchain() (*Toolchain, error) {
	embeddedOnce.Do(func() {
		embeddedTC, embeddedErr = setupEmbeddedToolchain()
	})
	return embeddedTC, embeddedErr
}

func setupEmbeddedToolchain() (*Toolchain, error) {
	dir, err := os.MkdirTemp("", "fuji-embedded-*")
	if err != nil {
		return nil, fmt.Errorf("embedded toolchain temp dir: %w", err)
	}
	platform := platformKey()
	files := map[string]string{
		llcBinaryName():     filepath.Join(dir, llcBinaryName()),
		lldBinaryName():     filepath.Join(dir, lldBinaryName()),
		"libfuji_runtime.a": filepath.Join(dir, "libfuji_runtime.a"),
	}
	for srcLeaf, dstPath := range files {
		// embed.FS paths always use forward slashes (see go:embed).
		srcPath := "bundled/" + platform + "/" + srcLeaf
		data, err := bundledReleaseFS.ReadFile(srcPath)
		if err != nil {
			return nil, fmt.Errorf("%w: read %s: %v", ErrIncompleteEmbeddedToolchain, srcPath, err)
		}
		mode := os.FileMode(0o644)
		if filepath.Ext(srcLeaf) != ".a" {
			mode = 0o755
		}
		if err := os.WriteFile(dstPath, data, mode); err != nil {
			return nil, fmt.Errorf("write embedded %s: %w", dstPath, err)
		}
	}

	llcPath := filepath.Join(dir, llcBinaryName())
	lldPath := filepath.Join(dir, lldBinaryName())
	libPath := filepath.Join(dir, "libfuji_runtime.a")

	tc := &Toolchain{
		LLC:        llcPath,
		LLD:        lldPath,
		RuntimeLib: libPath,
	}
	switch runtime.GOOS {
	case "linux":
		tc.LinkMode = LinkLLDGNU
		tc.Clang = ""
	case "darwin":
		tc.LinkMode = LinkLLDDarwin
		tc.Clang = ""
	case "windows":
		tc.LinkMode = LinkClang
		tc.Clang = findClangForEmbeddedLink(llcPath)
	default:
		tc.LinkMode = LinkClang
		tc.Clang = findClangForEmbeddedLink(llcPath)
	}
	return tc, nil
}
