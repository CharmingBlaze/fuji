//go:build !release

package fujembed

func Extract() (string, error) {
	return "", ErrDevelopmentBuild
}

func ClangPath() (string, error) {
	return "", ErrDevelopmentBuild
}

func LLCPath() (string, error) {
	return "", ErrDevelopmentBuild
}

func RuntimeLibPath() (string, error) {
	return "", ErrDevelopmentBuild
}

func LLDPathWindows() (string, error) {
	return "", ErrDevelopmentBuild
}
