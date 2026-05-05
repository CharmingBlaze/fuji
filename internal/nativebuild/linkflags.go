package nativebuild

import "runtime"

func defaultSystemLinkFlags() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"-lopengl32", "-lgdi32", "-lwinmm"}
	case "darwin":
		return []string{"-framework", "OpenGL", "-framework", "Cocoa", "-framework", "IOKit", "-framework", "CoreVideo"}
	case "linux":
		return []string{"-lGL", "-lm", "-lpthread", "-ldl", "-lrt", "-lX11"}
	default:
		return nil
	}
}
