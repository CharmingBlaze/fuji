package api

import (
	"path/filepath"

	"fuji/internal/nativebuild"
)

// BuildNativeHost compiles entryPath to a native executable (same pipeline as fuji build):
// Go codegen emits LLVM IR, llc lowers it to an object file, then Clang links that object with
// runtime/libfuji_runtime.a and headers under runtime/src (plus optional FUJI_NATIVE_SOURCES / FUJI_LINKFLAGS).
func BuildNativeHost(entryPath, overlay, output string, log func(string)) error {
	absEntry, err := filepath.Abs(entryPath)
	if err != nil {
		return err
	}
	overlays := map[string]string{}
	if overlay != "" {
		overlays[absEntry] = overlay
	}
	return nativebuild.BuildWithOverlays(entryPath, overlays, output, log)
}

// DefaultExeName returns the default native executable name for a .fuji path.
func DefaultExeName(source string) string {
	return nativebuild.DefaultExeName(source)
}
