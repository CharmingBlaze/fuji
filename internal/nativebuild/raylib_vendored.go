package nativebuild

import (
	"os"
	"path/filepath"
	"strings"
)

// vendoredRaylibStatic returns include dir and static archive path when the
// third_party/raylib_static/stage tree exists (from `make -C third_party/raylib_static`).
// Set FUJI_USE_VENDORED_RAYLIB=0 or false to skip even if stage exists.
func vendoredRaylibStatic(rootDir string) (includeDir, archive string, ok bool) {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("FUJI_USE_VENDORED_RAYLIB")))
	if v == "0" || v == "false" || v == "no" {
		return "", "", false
	}
	stage := filepath.Join(rootDir, "third_party", "raylib_static", "stage")
	inc := filepath.Join(stage, "include")
	candidates := []string{
		filepath.Join(stage, "lib", "libraylib.a"),
		filepath.Join(stage, "lib", "raylib.lib"),
	}
	for _, a := range candidates {
		if st, err := os.Stat(a); err == nil && !st.IsDir() {
			if fi, err := os.Stat(inc); err == nil && fi.IsDir() {
				return inc, a, true
			}
		}
	}
	return "", "", false
}
