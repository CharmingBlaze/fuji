package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"fuji/internal/diagnostic"
	"fuji/internal/fujihome"
	"fuji/internal/lexer"
)

// Parsed AST cache: absolute path → last seen mtime (nanos) + parsed program.
// Entries are invalidated when Stat mtime differs. Overlay sources are never cached.
type parseCacheEntry struct {
	modTimeNanos int64
	prog         *Program
}

var parseCacheMu sync.Mutex
var parseCache = make(map[string]parseCacheEntry)

func resetParseCache() {
	parseCacheMu.Lock()
	defer parseCacheMu.Unlock()
	parseCache = make(map[string]parseCacheEntry)
}

// LoadProgram parses the entry file and all its transitive imports.
func LoadProgram(entryPath string) (*ProgramBundle, error) {
	return LoadProgramWithOverlays(entryPath, nil)
}

// LoadProgramWithOverlays is like LoadProgram but allows providing source text for specific files.
func LoadProgramWithOverlays(entryPath string, overlays map[string]string) (*ProgramBundle, error) {
	absEntry, err := filepath.Abs(entryPath)
	if err != nil {
		return nil, err
	}

	bundle := &ProgramBundle{
		Modules: make(map[string]*Program),
	}

	visited := make(map[string]bool)
	if err := loadModule(absEntry, bundle, visited, overlays); err != nil {
		return nil, err
	}

	bundle.Entry = bundle.Modules[absEntry]
	return bundle, nil
}

func parseProgramSource(absPath string, src string) (*Program, error) {
	l := lexer.NewLexer(src, absPath)
	tokens, err := l.Tokenize()
	if err != nil {
		return nil, diagnostic.WrapLexer(absPath, src, err)
	}
	pr := NewParser(tokens)
	prog, err := pr.Parse()
	if err != nil {
		return nil, diagnostic.WrapParse(absPath, src, err)
	}
	return prog, nil
}

func loadModuleImports(modulePath string, prog *Program, bundle *ProgramBundle, visited map[string]bool, overlays map[string]string) error {
	var imports []string
	findImports(prog, &imports)
	for _, rel := range imports {
		abs, err := ResolveImportPath(modulePath, rel)
		if err != nil {
			return err
		}
		if err := loadModule(abs, bundle, visited, overlays); err != nil {
			return err
		}
	}
	return nil
}

func loadModule(path string, bundle *ProgramBundle, visited map[string]bool, overlays map[string]string) error {
	if visited[path] {
		if _, ok := bundle.Modules[path]; ok {
			return nil
		}
		return fmt.Errorf("import cycle detected: %s", path)
	}

	visited[path] = true
	defer delete(visited, path)

	var src string
	hasOverlayEntry := false
	if overlays != nil {
		if s, ok := overlays[path]; ok {
			hasOverlayEntry = true
			src = s
		}
	}

	if src == "" {
		if !hasOverlayEntry {
			parseCacheMu.Lock()
			e, cached := parseCache[path]
			parseCacheMu.Unlock()

			fi, statErr := os.Stat(path)
			if statErr == nil && cached && e.modTimeNanos == fi.ModTime().UnixNano() {
				bundle.Modules[path] = e.prog
				return loadModuleImports(path, e.prog, bundle, visited, overlays)
			}
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read %s: %w", path, err)
		}
		src = string(b)
	}

	prog, err := parseProgramSource(path, src)
	if err != nil {
		return err
	}
	bundle.Modules[path] = prog

	if !hasOverlayEntry {
		if fi, err := os.Stat(path); err == nil {
			parseCacheMu.Lock()
			parseCache[path] = parseCacheEntry{modTimeNanos: fi.ModTime().UnixNano(), prog: prog}
			parseCacheMu.Unlock()
		}
	}

	return loadModuleImports(path, prog, bundle, visited, overlays)
}

func findImports(node Node, imports *[]string) {
	switch n := node.(type) {
	case *Program:
		for _, d := range n.Declarations {
			findImports(d, imports)
		}
	case *IncludeDecl:
		rel := strings.Trim(n.Path.Lexeme, `"'`)
		*imports = append(*imports, rel)
	case *LetDecl:
		if n.Init != nil {
			findImports(n.Init, imports)
		}
	case *ImportExpr:
		*imports = append(*imports, strings.Trim(n.Path.Lexeme, "\""))
	case *ExpressionStmt:
		findImports(n.Expr, imports)
	case *BlockStmt:
		for _, s := range n.Declarations {
			findImports(s, imports)
		}
		// ... add other nodes as needed, or use a general walker.
	}
}

// pathListEnv reads a PATH-style environment variable (colon- or semicolon-separated).
func pathListEnv(key string) []string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return nil
	}
	return strings.Split(v, string(os.PathListSeparator))
}

// tryResolveModuleUnderRoots looks for moduleID (e.g. "array" or "raylib/core") under each root.
func tryResolveModuleUnderRoots(moduleID string, roots []string) (string, bool) {
	safe := filepath.FromSlash(moduleID)
	base := filepath.Base(safe)
	for _, root := range roots {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}
		candidates := []string{
			filepath.Join(root, safe+".fuji"),
			filepath.Join(root, safe, "index.fuji"),
			filepath.Join(root, safe, base+".fuji"),
		}
		for _, p := range candidates {
			fi, err := os.Stat(p)
			if err != nil || fi.IsDir() {
				continue
			}
			abs, err := filepath.Abs(p)
			if err != nil {
				continue
			}
			return abs, true
		}
	}
	return "", false
}

// ResolveImportPath resolves a relative or @module path relative to the importer.
func ResolveImportPath(importerPath, relPath string) (string, error) {
	if strings.HasPrefix(relPath, "@") {
		name := strings.ToLower(relPath[1:])

		// 1) Explicit search paths (override)
		if p, ok := tryResolveModuleUnderRoots(name, pathListEnv("FUJI_WRAPPERS")); ok {
			return p, nil
		}
		if p, ok := tryResolveModuleUnderRoots(name, pathListEnv("FUJI_PATH")); ok {
			return p, nil
		}

		// 2) Shipped layout next to fuji / stdlib/, wrappers/
		if inst, err := fujihome.InstallDir(); err == nil {
			bundled := []string{
				filepath.Join(inst, "stdlib"),
				filepath.Join(inst, "wrappers"),
				filepath.Join(inst, "lib"),
			}
			if p, ok := tryResolveModuleUnderRoots(name, bundled); ok {
				return p, nil
			}
		}

		// 3) Next to the importing file
		dir := filepath.Dir(importerPath)
		if p, ok := tryResolveModuleUnderRoots(name, []string{dir}); ok {
			return p, nil
		}

		return "", fmt.Errorf("could not resolve @ module %s", relPath)
	}

	return resolvePlainIncludePath(importerPath, relPath)
}

func resolvePlainIncludePath(importerPath, relPath string) (string, error) {
	if filepath.IsAbs(relPath) {
		return filepath.Clean(relPath), nil
	}

	tryFile := func(p string) (string, bool) {
		fi, err := os.Stat(p)
		if err != nil || fi.IsDir() {
			return "", false
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			return "", false
		}
		return abs, true
	}

	if p, ok := tryFile(filepath.Join(filepath.Dir(importerPath), relPath)); ok {
		return p, nil
	}

	for _, r := range pathListEnv("FUJI_PATH") {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		if p, ok := tryFile(filepath.Join(r, relPath)); ok {
			return p, nil
		}
	}

	if inst, err := fujihome.InstallDir(); err == nil {
		for _, base := range []string{
			filepath.Join(inst, "wrappers"),
			filepath.Join(inst, "stdlib"),
			inst,
		} {
			if p, ok := tryFile(filepath.Join(base, relPath)); ok {
				return p, nil
			}
		}
	}

	return filepath.Abs(filepath.Join(filepath.Dir(importerPath), relPath))
}

func BundleEntryPath(bundle *ProgramBundle) (string, error) {
	for path, prog := range bundle.Modules {
		if prog == bundle.Entry {
			return path, nil
		}
	}
	return "", fmt.Errorf("entry path not found in bundle")
}
