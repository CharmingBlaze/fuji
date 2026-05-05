package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fuji/internal/diagnostic"
	"fuji/internal/fujihome"
	"fuji/internal/lexer"
)

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

func loadModule(path string, bundle *ProgramBundle, visited map[string]bool, overlays map[string]string) error {
	if visited[path] {
		// Cycle detection could be more sophisticated, but for now we just return if already visiting.
		// Wait, if it's already in bundle.Modules, we're done.
		if _, ok := bundle.Modules[path]; ok {
			return nil
		}
		return fmt.Errorf("import cycle detected: %s", path)
	}

	visited[path] = true
	defer delete(visited, path)

	var src string
	if overlays != nil {
		if s, ok := overlays[path]; ok {
			src = s
		}
	}

	if src == "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read %s: %w", path, err)
		}
		src = string(b)
	}

	l := lexer.NewLexer(src)
	tokens, err := l.Tokenize()
	if err != nil {
		return diagnostic.WrapLexer(path, src, err)
	}

	p := NewParser(tokens)
	prog, err := p.Parse()
	if err != nil {
		return diagnostic.WrapParse(path, src, err)
	}

	bundle.Modules[path] = prog

	var imports []string
	findImports(prog, &imports)

	for _, rel := range imports {
		abs, err := ResolveImportPath(path, rel)
		if err != nil {
			return err
		}
		if err := loadModule(abs, bundle, visited, overlays); err != nil {
			return err
		}
	}

	return nil
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

	dir := filepath.Dir(importerPath)
	abs := filepath.Join(dir, relPath)
	return filepath.Abs(abs)
}

func BundleEntryPath(bundle *ProgramBundle) (string, error) {
	for path, prog := range bundle.Modules {
		if prog == bundle.Entry {
			return path, nil
		}
	}
	return "", fmt.Errorf("entry path not found in bundle")
}
