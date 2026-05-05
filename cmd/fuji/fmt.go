package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"fuji/internal/formatter"
)

var fmtSkipDirNames = map[string]bool{
	".git":                true,
	".FUJI_build":         true,
	"bin":                 true,
	"node_modules":        true,
	"vendor":              true,
	"_wg_legacy":          true,
	"_wg_out":             true,
	"raylib_full_wrapper": true,
	"test_output":         true,
	"test":                true,
}

func runFmtCmd(args []string) error {
	var check bool
	var paths []string
	for _, a := range args {
		switch a {
		case "--check":
			check = true
		default:
			paths = append(paths, a)
		}
	}
	if len(paths) == 0 {
		return fmt.Errorf("usage: fuji fmt [--check] <file.fuji> [more files...] [./...]\n       fuji fmt [--check] ./...")
	}

	files, err := expandFmtTargets(paths)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("fuji fmt: no .fuji files matched")
	}

	var checkFailed bool
	for _, path := range files {
		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		text := normalizeNL(string(raw))
		out, err := formatter.Format(text)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if check {
			if out != text {
				fmt.Fprintf(os.Stderr, "not formatted: %s\n", path)
				checkFailed = true
			}
			continue
		}
		if out == text {
			continue
		}
		if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
		fmt.Println(path)
	}
	if checkFailed {
		return fmt.Errorf("fuji fmt --check: one or more files need formatting")
	}
	return nil
}

func normalizeNL(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\r", "\n")
}

func expandFmtTargets(paths []string) ([]string, error) {
	seen := map[string]struct{}{}
	for _, p := range paths {
		switch p {
		case "./...", "...":
			list, err := collectFujiFiles(".")
			if err != nil {
				return nil, err
			}
			for _, f := range list {
				seen[f] = struct{}{}
			}
			continue
		}
		st, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		if st.IsDir() {
			list, err := collectFujiFiles(p)
			if err != nil {
				return nil, err
			}
			for _, f := range list {
				seen[f] = struct{}{}
			}
			continue
		}
		if !strings.HasSuffix(strings.ToLower(p), ".fuji") {
			return nil, fmt.Errorf("%s: expected a .fuji file or directory", p)
		}
		seen[p] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for f := range seen {
		out = append(out, f)
	}
	sort.Strings(out)
	return out, nil
}

func collectFujiFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root && fmtSkipDirNames[d.Name()] {
				return fs.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), ".fuji") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}
